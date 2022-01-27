import secrets
from datetime import datetime, timedelta

from django.conf import settings
from django.contrib.auth.hashers import check_password
from django.core.cache import cache
from django.db import transaction, IntegrityError
from django.http import Http404
from django.shortcuts import get_object_or_404
from django.utils.dateparse import parse_datetime
from rest_framework import status
from rest_framework.authtoken.models import Token
from rest_framework.decorators import api_view, permission_classes
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response
from rest_framework.views import APIView

from api.models import Tag, Hole, Floor, Report, User, Message, Division, PushToken, OldUserFavorites
from api.serializers import TagSerializer, HoleSerializer, FloorSerializer, ReportSerializer, MessageSerializer, \
    UserSerializer, DivisionSerializer, FloorGetSerializer, RegisterSerializer, EmailSerializer, BaseEmailSerializer, HoleCreateSerializer, \
    PushTokenSerializer, FloorUpdateSerializer
from api.signals import modified_by_admin, new_penalty, mention_to
from api.tasks import send_email
from utils.apis import find_mentions
from utils.auth import check_api_key, many_hashes
from utils.notification import send_notifications
from utils.permissions import OnlyAdminCanModify, OwnerOrAdminCanModify, NotSilentOrAdminCanPost, AdminOrReadOnly, \
    AdminOrPostOnly, OwenerOrAdminCanSee, AdminOnly


@api_view(["GET"])
def index(request):
    send_notifications.delay(request.user.id, 'hi')
    return Response({"message": "Hello world!"})


@api_view(["POST"])
def login(request):
    # TODO: Sanitize input needed?
    email = request.data.get("email")
    password = request.data.get("password")

    user = get_object_or_404(User, identifier=many_hashes(email))
    if check_password(password, user.password):
        token, created = Token.objects.get_or_create(user=user)
        return Response({"token": token.key, "message": "登录成功！"})
    else:
        return Response({"message": "用户名或密码错误！"}, 401)


@api_view(["GET"])
@permission_classes([IsAuthenticated])
def logout(request):
    request.auth.delete()
    Token.objects.create(user=request.user)
    return Response({"message": "登出成功"})


class VerifyApi(APIView):
    throttle_scope = 'email'

    @staticmethod
    def _set_verification_code(email: str) -> str:
        """
        设置验证码并返回
        """
        verification = secrets.randbelow(1000000)
        verification = str(verification).zfill(6)
        cache.set(email, verification, settings.VALIDATION_CODE_EXPIRE_TIME * 60)
        return verification

    def get(self, request, **kwargs):
        method = kwargs.get("method")

        serializer = EmailSerializer(data=request.query_params)
        serializer.is_valid(raise_exception=True)
        email = serializer.validated_data.get('email')
        uuid = serializer.validated_data.get('email')

        if method == "email":
            # 设置验证码并发送验证邮件
            verification = self._set_verification_code(email)
            base_content = (
                f'您的验证码是: {verification}\r\n'
                f'验证码的有效期为 {settings.VALIDATION_CODE_EXPIRE_TIME} 分钟\r\n'
                '如果您意外地收到了此邮件，请忽略它'
            )
            if not User.objects.filter(identifier=many_hashes(email)).exists():  # 用户不存在，注册邮件
                send_email.delay(
                    subject=f'{settings.SITE_NAME} 注册验证',
                    content=f'欢迎注册 {settings.SITE_NAME}，{base_content}',
                    receivers=[email],
                    uuid=uuid
                )
            else:  # 用户存在，重置密码
                send_email.delay(
                    subject=f'{settings.SITE_NAME} 重置密码',
                    content=f'您正在重置密码，{base_content}',
                    receivers=[email],
                    uuid=uuid
                )
            return Response({'message': '处理中'}, 202)
        elif method == "apikey":
            apikey = request.query_params.get("apikey")
            check_register = request.query_params.get("check_register")
            if not check_api_key(apikey):
                return Response({"message": "API Key 不正确！"}, 403)
            if not User.objects.filter(identifier=many_hashes(email)).exists():
                if check_register:
                    return Response({"message": "用户未注册！"}, 200)
                else:
                    verification = self._set_verification_code(email)
                    return Response({'message': '验证成功', 'code': verification}, 200)
            return Response({'message': '用户已注册'}, 409)
        else:
            return Response({}, 404)


class RegisterApi(APIView):
    def post(self, request):
        # TODO: Sanitize input needed?
        serializer = RegisterSerializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        user = serializer.save()
        token = Token.objects.get(user=user).key

        password = serializer.validated_data.get('password')
        email = serializer.validated_data.get('email')
        # 迁移用户收藏
        old_favorites = OldUserFavorites.objects.filter(uid=email[:11]).first()
        if old_favorites:
            try:
                user.favorites.set(old_favorites.favorites)
            except IntegrityError:
                pass

        # 发送密码邮件
        send_email.delay(
            subject=f'{settings.SITE_NAME} 密码存档',
            content=(
                f'您已成功注册{settings.SITE_NAME}，您选择了随机设置密码，密码如下：'
                f'\r\n\r\n{password}\r\n\r\n'
                '提示：服务器中仅存储加密后的密码，无须担心安全问题'
            ),
            receivers=[email]
        )
        return Response({'message': '注册成功', 'token': token}, 201)

    def put(self, request):
        email = request.data.get('email')
        user = get_object_or_404(User, identifier=many_hashes(email))
        serializer = RegisterSerializer(instance=user, data=request.data)
        serializer.is_valid(raise_exception=True)
        serializer.save()
        return Response({"message": "已重置密码"}, 200)


class EmailApi(APIView):
    throttle_scope = 'email'

    def post(self, request, **kwargs):
        serializer = BaseEmailSerializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        email = serializer.validated_data.get('email')
        uuid = serializer.validated_data.get('uuid')

        if kwargs.get('type') == 'password':
            password = request.data.get('password')
            if not password:
                return Response({'message': 'password 字段不存在'}, 400)
            send_email.delay(
                subject=f'{settings.SITE_NAME} 密码存档',
                content=(
                    f'您已成功注册{settings.SITE_NAME}，您选择了随机设置密码，密码如下：'
                    f'\r\n\r\n{password}\r\n\r\n'
                    '提示：服务器中仅存储加密后的密码，无须担心安全问题'
                ),
                receivers=[email],
                uuid=uuid
            )
            return Response({'message': '处理中'}, 202)
        else:
            raise Http404()


class DivisionsApi(APIView):
    permission_classes = [IsAuthenticated, AdminOrReadOnly]

    def get(self, request, **kwargs):
        division_id = kwargs.get('division_id')
        if division_id:
            query_set = get_object_or_404(Division, id=division_id)
        else:
            query_set = Division.objects.all()

        serializer = DivisionSerializer(query_set, many=not division_id, context={'user': request.user})
        return Response(serializer.data)

    @transaction.atomic
    def put(self, request, **kwargs):
        division_id = kwargs.get('division_id')
        division = get_object_or_404(Division, id=division_id)

        name = request.data.get('name')
        description = request.data.get('description')
        pinned = request.data.get('pinned')
        if name:
            division.name = name
        if description:
            division.description = description
        if pinned:
            division.pinned = pinned

        division.save()
        serializer = DivisionSerializer(division, context={'user': request.user})
        return Response(serializer.data)


class HolesApi(APIView):
    permission_classes = [IsAuthenticated, NotSilentOrAdminCanPost, OnlyAdminCanModify]

    def get(self, request, **kwargs):
        # 校验数据
        serializer = HoleSerializer(data=request.query_params)
        serializer.is_valid(raise_exception=True)
        length = serializer.validated_data.get('length')
        prefetch_length = serializer.validated_data.get('prefetch_length')
        start_time = serializer.validated_data.get('start_time')
        division_id = serializer.validated_data.get('division_id')

        # 获取单个
        hole_id = kwargs.get('hole_id')
        if hole_id:
            hole = get_object_or_404(Hole, pk=hole_id)
            # 缓存中增加主题帖的浏览量
            cached = cache.get('hole_views', {})
            view = cached.get(hole_id, 0)
            cached[hole_id] = view + 1
            cache.set('hole_views', cached)

            serializer = HoleSerializer(hole, context={
                "user": request.user,
                "prefetch_length": prefetch_length,
                'simple': False  # 不使用缓存
            })
            return Response(serializer.data)

        # 获取多个
        else:
            tag_name = request.query_params.get('tag')
            if tag_name:
                tag = get_object_or_404(Tag, name=tag_name)
                queryset = tag.hole_set.all()
            else:
                queryset = Hole.objects.all()

            queryset = queryset.order_by('-time_updated').filter(
                time_updated__lt=start_time,
                division_id=division_id
            )[:length]
            queryset = HoleSerializer.get_queryset(queryset)
            serializer = HoleSerializer(queryset, many=True, context={
                "user": request.user,
                "prefetch_length": prefetch_length,
                'simple': True  # 使用 SimpleFloorSerializer
            })
            return Response(serializer.data)

    @transaction.atomic
    def post(self, request):
        serializer = HoleCreateSerializer(data=request.data, context={'request_data': request.data, 'user': request.user})
        serializer.is_valid(raise_exception=True)
        # 检查权限
        division_id = serializer.validated_data.get('division_id')
        self.check_object_permissions(request, division_id)

        serializer.save()
        return Response({'message': '发表成功！', 'data': serializer.data}, 201)

    @transaction.atomic
    def put(self, request, **kwargs):
        hole_id = kwargs.get('hole_id')
        hole = get_object_or_404(Hole, pk=hole_id)

        serializer = HoleSerializer(hole, data=request.data)
        serializer.is_valid(raise_exception=True)
        hole = serializer.save()

        serializer = HoleSerializer(hole, context={"user": request.user})
        return Response(serializer.data)

    def delete(self, request, **kwargs):
        # 主题帖不能删除
        return Response(None, status=status.HTTP_405_METHOD_NOT_ALLOWED)


class FloorsApi(APIView):
    permission_classes = [IsAuthenticated, NotSilentOrAdminCanPost, OwnerOrAdminCanModify]

    def get(self, request, **kwargs):
        # 获取单个
        floor_id = kwargs.get('floor_id')
        if floor_id:
            floor = get_object_or_404(Floor, pk=floor_id)
            serializer = FloorSerializer(floor, context={"user": request.user})
            return Response(serializer.data)
        # 获取多个
        serializer = FloorGetSerializer(data=request.query_params)
        serializer.is_valid(raise_exception=True)
        hole_id = serializer.validated_data.get('hole_id')
        search = serializer.validated_data.get('s')
        start_floor = serializer.validated_data.get('start_floor')
        length = serializer.validated_data.get('length')
        reverse = serializer.validated_data.get('reverse')

        if search:  # 搜索
            query_set = Floor.objects.filter(shadow_text__icontains=search)
            if not reverse:  # 搜索默认降序，reverse 反转
                query_set = query_set.order_by('-pk')
        else:  # 主题帖下
            query_set = Floor.objects.filter(hole_id=hole_id)
            if reverse:  # 主题帖默认升序，reverse 反转
                query_set = query_set.order_by('-pk')

        if length == 0:
            query_set = query_set[start_floor:]
        else:
            query_set = query_set[start_floor: start_floor + length]

        query_set = FloorSerializer.get_queryset(query_set)
        serializer = FloorSerializer(query_set, many=True, context={"user": request.user})
        return Response(serializer.data)

    @transaction.atomic
    def post(self, request):
        hole_id = request.data.get('hole_id')
        hole = get_object_or_404(Hole, pk=hole_id)
        self.check_object_permissions(request, hole.division_id)
        serializer = FloorSerializer(data=request.data, context={'user': request.user, 'hole': hole})
        serializer.is_valid(raise_exception=True)
        serializer.save()
        return Response({'message': '发表成功！', 'data': serializer.data}, 201)

    @transaction.atomic
    def put(self, request, **kwargs):
        floor_id = kwargs.get('floor_id')
        floor = get_object_or_404(Floor, pk=floor_id)
        serializer = FloorUpdateSerializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        data = serializer.validated_data

        # 不检查权限
        like = data.pop('like', '')
        if like:
            if like == 'add' and request.user.pk not in floor.like_data:
                floor.like_data.append(request.user.pk)
            elif like == 'cancel' and request.user.pk in floor.like_data:
                floor.like_data.remove(request.user.pk)
            else:
                pass
            floor.like = len(floor.like_data)

        # 属主或管理员
        if data:
            self.check_object_permissions(request, floor)
        content = data.pop('content', '')
        if content:
            floor.history.append({
                'content': floor.content,
                'altered_by': request.user.id,
                'altered_time': datetime.now(settings.TIMEZONE).isoformat()
            })
            floor.content = content
            mentions = find_mentions(content)
            floor.mention.set(mentions)
            mention_to.send(sender=Floor, instance=floor, mentioned=mentions)
        floor.fold = data.pop('fold', floor.fold)

        # 仅管理员
        if data and not request.user.is_admin:
            return Response(None, 403)
        floor.anonyname = data.pop('anonyname', floor.anonyname)
        floor.special_tag = data.pop('special_tag', floor.special_tag)

        floor.save()

        if request.user.is_admin and floor.user != request.user:
            modified_by_admin.send(sender=Floor, instance=floor)
        return Response(FloorSerializer(floor, context={'user': request.user}).data)

    @transaction.atomic
    def delete(self, request, **kwargs):
        floor_id = kwargs.get('floor_id')
        delete_reason = request.data.get('delete_reason')
        floor = get_object_or_404(Floor, pk=floor_id)
        self.check_object_permissions(request, floor)
        floor.history.append({
            'content': floor.content,
            'altered_by': request.user.pk,
            'altered_time': datetime.now(settings.TIMEZONE).isoformat()
        })
        if request.user == floor.user:  # 作者删除
            floor.content = '该内容已被作者删除'
        else:  # 管理员删除
            floor.content = delete_reason if delete_reason else '该内容因违反社区规范被删除'
        floor.deleted = True
        floor.save()
        serializer = FloorSerializer(floor, context={"user": request.user})
        return Response(serializer.data, 200)


class TagsApi(APIView):
    permission_classes = [IsAuthenticated, AdminOrReadOnly]

    def get(self, request, **kwargs):
        # 获取单个
        tag_id = kwargs.get('tag_id')
        if tag_id:
            tag = get_object_or_404(Tag, pk=tag_id)
            serializer = TagSerializer(tag)
            return Response(serializer.data)
        # 获取列表
        search = request.query_params.get('s')
        query_set = Tag.objects.order_by('-temperature')
        if search:
            query_set = query_set.filter(name__icontains=search)
        serializer = TagSerializer(query_set, many=True)
        return Response(serializer.data)

    def post(self, request):
        serializer = TagSerializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        serializer.save()
        return Response(serializer.data, 201)

    def put(self, request, **kwargs):
        tag_id = kwargs.get('tag_id')
        tag = get_object_or_404(Tag, pk=tag_id)

        serializer = TagSerializer(instance=tag, data=request.data)
        serializer.is_valid(raise_exception=True)
        serializer.save()

        return Response(serializer.data)

    def delete(self, request, **kwargs):
        tag_id = kwargs.get('tag_id')
        tag = get_object_or_404(Tag, pk=tag_id)
        tag.delete()
        return Response(None, 204)


class FavoritesApi(APIView):
    permission_classes = [IsAuthenticated]

    def get(self, request):
        query_set = request.user.favorites.all()
        serializer = HoleSerializer(query_set, many=True, context={"user": request.user, 'simple': True})
        return Response(serializer.data)

    def post(self, request):
        hole_id = request.data.get('hole_id')
        hole = get_object_or_404(Hole, pk=hole_id)
        request.user.favorites.add(hole)
        return Response({'message': '收藏成功', 'data': request.user.favorites.values_list('id', flat=True)}, 201)

    def put(self, request):
        hole_ids = request.data.get('hole_ids')
        holes = Hole.objects.filter(pk__in=hole_ids)
        request.user.favorites.set(holes)
        return Response({'message': '修改成功', 'data': request.user.favorites.values_list('id', flat=True)}, 200)

    def delete(self, request):
        hole_id = request.data.get('hole_id')
        hole = get_object_or_404(Hole, pk=hole_id)
        request.user.favorites.remove(hole)
        return Response({'message': '删除成功', 'data': request.user.favorites.values_list('id', flat=True)}, 200)


class ReportsApi(APIView):
    permission_classes = [IsAuthenticated, AdminOrPostOnly]

    def post(self, request):
        floor_id = request.data.get('floor_id')
        reason = request.data.get('reason')
        floor = get_object_or_404(Floor, pk=floor_id)
        if not reason or not reason.strip():
            return Response({'message': '举报原因不能为空'}, 400)
        report = Report.objects.create(hole_id=floor.hole_id, floor_id=floor_id, reason=reason)
        serializer = ReportSerializer(report)
        return Response(serializer.data, 201)

    def get(self, request, **kwargs):
        # 获取单个
        report_id = kwargs.get('report_id')
        if report_id:
            report = get_object_or_404(Report, pk=report_id)
            serializer = ReportSerializer(report)
            return Response(serializer.data)
        # 获取多个
        category = request.query_params.get('category', default='not_dealed')
        if category == 'not_dealed':
            queryset = Report.objects.filter(dealed=False)
        elif category == 'dealed':
            queryset = Report.objects.filter(dealed=True)
        elif category == 'all':
            queryset = Report.objects.all()
        else:
            return Response({'message': 'category 参数不正确'})
        queryset = queryset.order_by('-time_created')
        serializer = ReportSerializer(queryset, many=True)
        return Response(serializer.data)

    def delete(self, request, **kwargs):
        report_id = kwargs.get('report_id')
        report = get_object_or_404(Report, pk=report_id)
        floor = report.floor

        if request.data.get('not_deal'):
            pass
        if request.data.get('fold'):
            floor.fold = request.data.get('fold')
        if request.data.get('delete'):
            delete_reason = request.data.get('delete')
            floor.history.append({
                'content': floor.content,
                'altered_by': request.user.pk,
                'altered_time': datetime.now(settings.TIMEZONE).isoformat()
            })
            floor.content = delete_reason
            floor.deleted = True
        if request.data.get('silent'):
            permission = floor.user.permission
            current_time_str = permission['silent'].get(str(floor.hole.division_id), '1970-01-01T00:00:00+00:00')
            current_time = parse_datetime(current_time_str)
            expected_time = datetime.now(settings.TIMEZONE) + timedelta(days=request.data.get('silent'))
            permission['silent'][str(floor.hole.division_id)] = max(current_time, expected_time).isoformat()
            floor.user.save()

        floor.save()
        report.dealed_by = request.user
        report.dealed = True
        report.save()
        return Response({'message': '举报处理成功'}, 200)

    def put(self, request):
        pass


class ImagesApi(APIView):
    permission_classes = [IsAuthenticated]

    def post(self, request):
        return Response({'message': '该 API 已弃用，请调用 websocket API 以上传图片'}, 405)


class MessagesApi(APIView):
    permission_classes = [IsAuthenticated, OnlyAdminCanModify, OwenerOrAdminCanSee]

    def post(self, request):
        floor = get_object_or_404(Floor, pk=request.data.get('to'))
        to_id = floor.user.pk

        if request.data.get('share_email'):
            code = 'share_email'
            message = f'用户看到了你发布的帖子\n{str(floor)}\n希望与你取得联系，TA的邮箱为：{request.user.email}'
        elif request.data.get('message'):
            code = 'message'
            message = request.data.get('message').strip()
            if not request.user.is_admin:
                return Response(None, 403)
            if not message:
                return Response({'message': 'message不能为空'}, 400)
        else:
            return Response(None, 400)

        Message.objects.create(user_id=to_id, message=message, code=code)
        return Response({'message': f'已发送通知，内容为：{message}'}, 201)

    def get(self, request, **kwargs):
        not_read = request.query_params.get('not_read', False)
        start_time = request.query_params.get('start_time')
        message_id = kwargs.get('message_id')

        # 获取单个
        if message_id:
            message = get_object_or_404(Message, pk=message_id)
            self.check_object_permissions(request, message)
            serializer = MessageSerializer(message)
            return Response(serializer.data)
        # 获取多个
        else:
            query_set = Message.objects.filter(user=request.user).order_by('-pk')
            if not_read:
                query_set = query_set.filter(has_read=False)
            if start_time:
                query_set = query_set.filter(time_created__lt=start_time)
            length = settings.FLOOR_PREFETCH_LENGTH
            serializer = MessageSerializer(query_set[:length], many=True)
            return Response(serializer.data)

    def put(self, request, **kwargs):
        message_id = kwargs.get('message_id')
        if message_id:
            message = get_object_or_404(Message, pk=message_id)
            content = request.data.get('message')
            has_read = request.data.get('has_read')
            code = request.data.get('code')
            data = request.data.get('data')

            if content:
                message.message = content.strip()
            if has_read:
                message.has_read = has_read
            if code:
                message.code = code
            if data:
                message.data = data

            message.save()
            serializer = MessageSerializer(message)
            return Response(serializer.data)
        else:
            clear_all = request.data.get('clear_all', False)
            if clear_all:
                Message.objects.filter(user=request.user).update(has_read=True)
                return Response(None, 200)
        return Response({'message': '需要指定操作'}, 400)

    def delete(self, request):
        pass


class UsersApi(APIView):
    permission_classes = [IsAuthenticated, OwnerOrAdminCanModify, OwenerOrAdminCanSee]

    def get(self, request, **kwargs):
        user_id = kwargs.get('user_id')
        if user_id:
            user = get_object_or_404(User, pk=user_id)
            self.check_object_permissions(request, user)
        else:
            user = request.user
        serializer = UserSerializer(user)
        return Response(serializer.data)

    def put(self, request, **kwargs):
        user_id = kwargs.get('user_id')
        if user_id:
            user = get_object_or_404(User, pk=user_id)
            self.check_object_permissions(request, user)
        else:
            user = request.user

        serializer = UserSerializer(data=request.data)
        serializer.is_valid()
        favorites = serializer.validated_data.get('favorites')
        config = serializer.validated_data.get('config')
        permission = serializer.validated_data.get('permission')

        if permission and request.user.is_admin:
            if user.permission == permission:
                pass  # 避免没有更改权限时发出信号
            else:
                user.permission = permission
                user.save(update_fields=['permission'])  # 发送权限被更改的信号
        if favorites:
            user.favorites.set(favorites)
        if config:
            user.config = config
        user.save()

        serializer = UserSerializer(user)
        return Response(serializer.data)


class PushTokensAPI(APIView):
    permission_classes = [IsAuthenticated]

    def get(self, request):
        if not request.user.is_admin:
            return Response(None, 403)
        if request.query_params.get('user_id'):
            user = get_object_or_404(User, pk=request.query_params.get('user_id'))
        else:
            user = request.user
        tokens = PushToken.objects.filter(user=user)
        service = request.query_params.get('service')
        if service:
            tokens = PushToken.objects.filter(service=service)
        return Response(PushTokenSerializer(tokens, many=True).data)

    def put(self, request):
        device_id = request.data.get('device_id', '')
        service = request.data.get('service', '')
        token = request.data.get('token', '')
        push_token = PushToken.objects.filter(device_id=device_id, user=request.user).first()
        if not push_token:
            push_token = PushToken.objects.create(device_id=device_id, service=service, token=token, user=request.user)
            code = 201
        else:
            push_token.token = token or push_token.token
            push_token.service = service or push_token.service
            push_token.save()
            code = 200
        serializer = PushTokenSerializer(push_token)
        return Response(serializer.data, code)

    def delete(self, request):
        device_id = request.data.get('device_id', '')
        PushToken.objects.filter(user=request.user, device_id=device_id).delete()
        return Response(None, 204)


class PenaltyApi(APIView):
    permission_classes = [AdminOnly]

    def post(self, request, **kwargs):
        self.check_object_permissions(request, request.user)
        floor = get_object_or_404(Floor, pk=kwargs.get('floor_id'))
        user = floor.user

        try:
            penalty_level = int(request.data.get('penalty_level'))
            division_id = request.data.get('division_id')
        except (ValueError, TypeError):
            return Response(status=status.HTTP_400_BAD_REQUEST)
        if penalty_level > 0:
            penalty_multiplier = {
                1: 1,
                2: 5,
                3: 999
            }.get(penalty_level, 1)

            offense_count = user.permission.get('offense_count', 0)
            offense_count += 1
            user.permission['offense_count'] = offense_count

            new_penalty_date = datetime.now(settings.TIMEZONE) + timedelta(days=offense_count * penalty_multiplier)
            user.permission['silent'][str(division_id)] = new_penalty_date.isoformat()

            new_penalty.send(sender=Floor, instance=floor, penalty=(penalty_level, new_penalty_date, division_id))

        user.save(update_fields=['permission'])
        serializer = UserSerializer(user)
        return Response(serializer.data)
