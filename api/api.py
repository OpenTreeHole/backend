import base64
import random
import secrets
import uuid
from datetime import datetime, timezone, timedelta

import magic
from django.conf import settings
from django.contrib.auth.hashers import check_password
from django.contrib.auth.password_validation import validate_password
from django.core.cache import cache
from django.core.exceptions import ValidationError
from django.db import transaction
from django.db.models import F
from django.shortcuts import get_object_or_404
from django.utils.dateparse import parse_datetime
from rest_framework import status
from rest_framework.authtoken.models import Token
from rest_framework.decorators import api_view, permission_classes
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response
from rest_framework.views import APIView

from api.models import Tag, Hole, Floor, Report, User, Message, Division
from api.notification import MessageSender
from api.permissions import OnlyAdminCanModify, OwnerOrAdminCanModify, NotSilentOrAdminCanPost, AdminOrReadOnly, \
    AdminOrPostOnly, OwenerOrAdminCanSee, AdminOnly
from api.serializers import TagSerializer, HoleSerializer, FloorSerializer, ReportSerializer, MessageSerializer, \
    UserSerializer, DivisionSerializer, MentionSerializer, FloorGetSerializer
from api.signals import modified_by_admin
from api.tasks import mail, post_image_to_github
# 发送 csrf 令牌
# from django.views.decorators.csrf import ensure_csrf_cookie
# @method_decorator(ensure_csrf_cookie)
from api.utils import encrypt_email, check_api_key


@api_view(["GET"])
def index(request):
    MessageSender(request.user, {'message': 'hi'}).commit()
    return Response({"message": "Hello world!"})


@api_view(["POST"])
def login(request):
    # TODO: Sanitize input needed?
    email = request.data.get("email")
    password = request.data.get("password")

    user = get_object_or_404(User, email=encrypt_email(email))
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

    def get(self, request, **kwargs):
        method = kwargs.get("method")

        if method == "email":
            email = request.query_params.get("email")
            domain = email[email.find("@") + 1:]
            # 检查邮箱是否在白名单内
            if domain not in settings.EMAIL_WHITELIST:
                return Response({"message": "邮箱不在白名单内！"}, 400)
            # 仅当用户未注册时发送邮件
            if not User.objects.filter(email=encrypt_email(email)).exists():
                # 设置验证码并发送验证邮件
                verification = secrets.randbelow(1000000)
                cache.set(email, verification, settings.VALIDATION_CODE_EXPIRE_TIME * 60)
                mail.delay(
                    subject=f'{settings.SITE_NAME} 注册验证',
                    content=f'欢迎注册 {settings.SITE_NAME}，您的验证码是: {str(verification).zfill(6)}\r\n验证码的有效期为 {settings.VALIDATION_CODE_EXPIRE_TIME} 分钟\r\n如果您意外地收到了此邮件，请忽略它',
                    receivers=[email]
                )
            return Response({'message': '如果您未注册，则会收到包含验证码的验证邮件，请查收；如果您已注册，则不会收到验证邮件，请找回密码。'})
        elif method == "apikey":
            apikey = request.query_params.get("apikey")
            check_register = request.query_params.get("check_register")
            if not check_api_key(apikey):
                return Response({"message": "API Key 不正确！"}, 403)
            email = request.query_params.get("email")
            domain = email[email.find("@") + 1:]
            if domain not in settings.EMAIL_WHITELIST:
                return Response({"message": "邮箱不在白名单内！"}, 400)
            if not User.objects.filter(email=encrypt_email(email)).exists():
                if check_register:
                    return Response({"message": "用户未注册！"}, 200)
                else:
                    verification = secrets.randbelow(1000000)
                    cache.set(email, verification, settings.VALIDATION_CODE_EXPIRE_TIME * 60)
                    return Response({'message': '验证成功', 'code': str(verification).zfill(6)}, 200)
            return Response({'message': '用户已注册'}, 409)
        else:
            return Response({}, 404)


class RegisterApi(APIView):
    def post(self, request):
        # TODO: Sanitize input needed?
        email = request.data.get("email")
        password = request.data.get("password")
        verification = request.data.get("verification")
        if not verification:
            return Response({"message": "验证码不能为空！"}, 400)
        if not cache.get(email) or not cache.get(email) == verification:
            return Response({"message": "注册校验未通过！"}, 400)
        # 校验密码可用性
        try:
            validate_password(password)
        except ValidationError as e:
            return Response({'message': '\n'.join(e)}, 400)
        # 校验用户名是否已存在
        if User.objects.filter(email=email).exists():
            return Response({"message": "该用户已注册！如果忘记密码，请使用忘记密码功能找回。"}, 400)

        user = User.objects.create_user(email=encrypt_email(email), password=password)
        token = Token.objects.get(user=user).key
        # RegisteredEmail.objects.create(email_cleartext=email).save()
        return Response({'message': '注册成功', 'token': token}, 201)

    def put(self, request):
        email = request.data.get("email")
        password = request.data.get("password")
        verification = request.data.get("verification")

        # 校验验证码
        if not verification:
            return Response({"message": "验证码不能为空！"}, 400)
        if not cache.get(email) or not cache.get(email) == verification:
            return Response({"message": "注册校验未通过！"}, 400)
        # 校验密码可用性
        try:
            validate_password(password)
        except ValidationError as e:
            return Response({'message': '\n'.join(e)}, 400)
        # 校验用户名是否不存在
        users = User.objects.filter(email=encrypt_email(email))
        if not users:
            return Response({"message": "该用户不存在"}, 400)
        user = users[0]

        user.set_password(password)
        user.save()
        return Response({"message": "已重置密码"}, 200)


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


def add_a_floor(request, hole, category):
    """
    增加一条回复帖
    Args:
        request:
        hole:       hole对象
        category:   指定返回值为 floor 或 hole

    Returns:        floor or hole

    """
    # 校验 content
    serializer = FloorSerializer(data=request.data)
    serializer.is_valid(raise_exception=True)
    content = serializer.validated_data.get('content')
    # 校验 mention
    mention_serializer = MentionSerializer(data=request.data)
    mention_serializer.is_valid(raise_exception=True)
    mention = mention_serializer.validated_data.get('mention', [])

    # 获取匿名信息，如没有则随机选取一个，并判断有无重复
    anonyname = hole.mapping.get(str(request.user.pk))  # 存在数据库中的字典里的数据类型都是 string
    if not anonyname:
        while True:
            anonyname = random.choice(settings.NAME_LIST)
            if anonyname in hole.mapping.values():
                pass
            else:
                hole.mapping[request.user.pk] = anonyname
                break
    hole.save()

    # 创建 floor 并增加 hole 的楼层数
    floor = Floor.objects.create(hole=hole, content=content, anonyname=anonyname, user=request.user)
    floor.mention.set(mention)
    return hole if category == 'hole' else floor


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
            Hole.objects.filter(pk=hole_id).update(view=F('view') + 1)  # 增加主题帖的浏览量
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
        serializer = HoleSerializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        tags = serializer.validated_data.get('tags')
        division_id = serializer.validated_data.get('division_id')
        self.check_object_permissions(request, division_id)

        # 在添加 tag 前要保存 hole，否则没有id
        hole = Hole.objects.create(division_id=division_id)

        # 创建 tag 并添加至 hole
        for tag in tags:
            try:
                tag = Tag.objects.get(name=tag['name'])
            except Tag.DoesNotExist:
                if 'color' not in tag:
                    tag.color = random.choice(settings.TAG_COLORS)
                tag = Tag.objects.create(name=tag['name'], color=tag['color'])
            hole.tags.add(tag)

        hole = add_a_floor(request, hole, category='hole')
        serializer = HoleSerializer(hole, context={"user": request.user})
        return Response({'message': '发表成功！', 'data': serializer.data}, 201)

    def put(self, request, **kwargs):
        hole_id = kwargs.get('hole_id')
        hole = get_object_or_404(Hole, pk=hole_id)
        serializer = HoleSerializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        tags = serializer.validated_data.get('tags')
        view = serializer.validated_data.get('view')

        if tags:
            hole.tags.clear()
            for tag in tags:
                try:
                    tag = Tag.objects.get(name=tag['name'])
                except Tag.DoesNotExist:
                    if 'color' not in tag:
                        tag.color = random.choice(settings.TAG_COLORS)
                    tag = Tag.objects.create(name=tag['name'], color=tag['color'])
                hole.tags.add(tag)
        if view:
            hole.view = view

        hole.save()
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

        if search:  # 搜索
            query_set = Floor.objects.filter(shadow_text__icontains=search).order_by('-pk')
        else:  # 主题帖下
            query_set = Floor.objects.filter(hole_id=hole_id)

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
        serializer = FloorSerializer(add_a_floor(request, hole, category='floor'), context={"user": request.user})
        return Response({'message': '发表成功！', 'data': serializer.data}, 201)

    def put(self, request, **kwargs):
        floor_id = kwargs.get('floor_id')
        content = request.data.get('content')
        like = request.data.get('like')
        fold = request.data.get('fold')
        mention = request.data.get('mention')
        anonyname = request.data.get('anonyname')
        floor = get_object_or_404(Floor, pk=floor_id)
        self.check_object_permissions(request, floor)

        if content and content.strip():
            floor.history.append({
                'content': floor.content,
                'altered_by': request.user.pk,
                'altered_time': datetime.now(timezone.utc).isoformat()
            })
            floor.content = content
        if like:
            if like == 'add' and request.user.pk not in floor.like_data:
                floor.like_data.append(request.user.pk)
            elif like == 'cancel' and request.user.pk in floor.like_data:
                floor.like_data.remove(request.user.pk)
            else:
                pass
            floor.like = len(floor.like_data)
        if mention:
            floor.mention.set(mention)
        if fold:
            floor.fold = fold
        if anonyname and request.user.is_admin:
            floor.anonyname = anonyname

        floor.save()
        if request.user.is_admin and floor.user != request.user:
            modified_by_admin.send(sender=Floor, instance=floor)
        serializer = FloorSerializer(floor, context={"user": request.user})
        return Response(serializer.data)

    def delete(self, request, **kwargs):
        floor_id = kwargs.get('floor_id')
        delete_reason = request.data.get('delete_reason')
        floor = get_object_or_404(Floor, pk=floor_id)
        self.check_object_permissions(request, floor)
        floor.history.append({
            'content': floor.content,
            'altered_by': request.user.pk,
            'altered_time': datetime.now(timezone.utc).isoformat()
        })
        if request.user == floor.user:  # 作者删除
            floor.content = '该内容已被作者删除'
        else:  # 管理员删除
            floor.content = delete_reason if delete_reason else '该内容因违反社区规范被删除'
        floor.deleted = True
        floor.save()
        serializer = FloorSerializer(floor, context={"user": request.user})
        return Response(serializer.data, 204)


class TagsApi(APIView):
    permission_classes = [IsAuthenticated, AdminOrReadOnly]

    def get(self, request):
        search = request.query_params.get('s')
        query_set = Tag.objects.order_by('-temperature')
        if search:
            query_set = query_set.filter(name__icontains=search)
        serializer = TagSerializer(query_set, many=True)
        return Response(serializer.data)

    def post(self, request):
        serializer = TagSerializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        name = serializer.validated_data.get('name')
        tag = Tag.objects.create(name=name, temperature=0)
        serializer = TagSerializer(tag)
        return Response(serializer.data, 201)

    def put(self, request, **kwargs):
        serializer = TagSerializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        name = serializer.validated_data.get('name')
        temperature = serializer.validated_data.get('temperature')
        tag_id = kwargs.get('tag_id')
        tag = get_object_or_404(Tag, pk=tag_id)
        if name:
            tag.name = name
        if temperature:
            tag.temperature = temperature
        serializer = TagSerializer(tag)
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
        serializer = ReportSerializer(queryset, many=True)
        return Response(serializer.data)

    def delete(self, request, **kwargs):
        report_id = kwargs.get('report_id')
        report = get_object_or_404(Report, pk=report_id)
        floor = report.floor
        deal = request.data.get('deal')

        if deal.get('not_deal'):
            pass
        if deal.get('fold'):
            floor.fold = deal.get('fold')
        if deal.get('delete'):
            delete_reason = deal.get('delete')
            floor.history.append({
                'content': floor.content,
                'altered_by': request.user.pk,
                'altered_time': datetime.now(timezone.utc).isoformat()
            })
            floor.content = delete_reason
            floor.deleted = True
        if deal.get('silent'):
            permission = floor.user.permission
            current_time_str = permission['silent'].get(str(floor.hole.division_id), '1970-01-01T00:00:00+00:00')
            current_time = parse_datetime(current_time_str)
            expected_time = datetime.now(timezone.utc) + timedelta(days=deal.get('silent'))
            permission['silent'][str(floor.hole.division_id)] = max(current_time, expected_time).isoformat()
            floor.user.save()

        floor.save()
        report.dealed_by = request.user
        report.dealed = True
        report.save()
        return Response({'message': '举报处理成功'}, 204)

    def put(self, request):
        pass


class ImagesApi(APIView):
    permission_classes = [IsAuthenticated]
    throttle_scope = 'upload'

    def post(self, request):
        # 校验图片
        image = request.data.get('image')
        if not image:
            return Response({'message': '内容不能为空'}, 400)
        if image.size > settings.MAX_IMAGE_SIZE * 1024 * 1024:
            return Response({'message': f'图片大小不能超过 {settings.MAX_IMAGE_SIZE} MB'}, 400)
        mime = magic.from_buffer(image.read(min([image.size, 2048])), mime=True)
        image.seek(0)
        if mime.split('/')[0] != 'image':
            return Response({'message': '请上传图片格式'}, 400)

        # 准备数据
        date_str = datetime.now().strftime('%Y-%m-%d')
        uid = uuid.uuid1()
        file_type = mime.split('/')[1]
        upload_url = f'https://api.github.com/repos/{settings.GITHUB_OWENER}/{settings.GITHUB_REPO}/contents/{date_str}/{uid}.{file_type}'
        headers = {
            'Authorization': f'token {settings.GITHUB_TOKEN}'
        }
        body = {
            'content': base64.b64encode(image.read()).decode('utf-8'),
            'message': f'upload image by user {request.user.pk}',
            'branch': settings.GITHUB_BRANCH,
        }
        post_image_to_github.delay(url=upload_url, headers=headers, body=body)

        result_url = f'https://cdn.jsdelivr.net/gh/{settings.GITHUB_OWENER}/{settings.GITHUB_REPO}@{settings.GITHUB_BRANCH}/{date_str}/{uid}.{file_type}'
        return Response({'url': result_url, 'message': '图片已上传'}, 202)


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
                query_set = query_set.filter(time_created__gt=start_time)

            serializer = MessageSerializer(query_set, many=True)
            return Response(serializer.data)

    def put(self, request, **kwargs):
        message_id = kwargs.get('message_id')
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

    def post(self, request, **kwargs):
        # This is (currently) used for Push notification token registration
        user_id = kwargs.get('user_id')
        if user_id:
            user = get_object_or_404(User, pk=user_id)
            self.check_object_permissions(request, user)
        else:
            user = request.user

        if ('service' not in request.data) or (request.data['service'] != 'apns' and request.data['service'] != 'mipush') or ('token' not in request.data) or (
                'device_id' not in request.data):
            return Response(status=status.HTTP_400_BAD_REQUEST)

        service = request.data['service']  # 'apns' or 'mipush'
        token = request.data['token']
        device_id = request.data['device_id']

        user.push_notification_tokens[service].update({device_id: token})
        user.save(update_fields=['push_notification_tokens'])

        return Response(status=status.HTTP_200_OK)


class PenaltyApi(APIView):
    permission_classes = [AdminOnly]

    def post(self, request, **kwargs):
        user_id = kwargs.get('user_id')
        user = get_object_or_404(User, pk=user_id)
        self.check_object_permissions(request, user)

        try:
            penalty_level = int(request.data.get('penalty_level'))
            division_id = request.data.get('division_id')
        except (ValueError, TypeError):
            return Response(status=status.HTTP_400_BAD_REQUEST)
        if penalty_level > 0:
            user.permission['offense_count'] += 1
            penalty_multiplier = 1
            if penalty_level == 2:
                penalty_multiplier = 5
            elif penalty_level == 3:
                penalty_multiplier = 999
            new_penalty_date = datetime.now() + timedelta(days=int(user.permission['offense_count']) * penalty_multiplier)
            user.permission['silent'].update({division_id: new_penalty_date.isoformat()})

        user.save(update_fields=['permission'])
        serializer = UserSerializer(user)
        return Response(serializer.data)
