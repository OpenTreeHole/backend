import base64
import random
import uuid
from datetime import datetime, timezone, timedelta

import magic
from django.conf import settings
from django.contrib.auth.hashers import check_password
from django.contrib.auth.models import User
from django.contrib.auth.password_validation import validate_password
from django.core.cache import cache
from django.core.exceptions import ValidationError
from django.db.models import F
from django.shortcuts import get_object_or_404
from django.utils.dateparse import parse_datetime
from rest_framework.authtoken.models import Token
from rest_framework.decorators import api_view
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response
from rest_framework.views import APIView

from api.models import Tag, Hole, Floor, Report, Profile
from api.permissions import OnlyAdminCanModify, OwnerOrAdminCanModify, NotSilentOrAdminCanPost, AdminOrReadOnly, AdminOrPostOnly, is_permitted
from api.serializers import TagSerializer, HoleSerializer, FloorSerializer, ReportSerializer
from api.tasks import hello_world, mail, post_image_to_github, send_message
from api.utils import to_shadow_text


# 发送 csrf 令牌
# from django.views.decorators.csrf import ensure_csrf_cookie
# @method_decorator(ensure_csrf_cookie)


@api_view(["GET"])
def index(request):
    hello_world.delay()
    return Response({"message": "Hello world!"})


@api_view(["POST"])
def login(request):
    username = request.data.get("email")
    password = request.data.get("password")
    user = get_object_or_404(User, username=username)
    if check_password(password, user.password):
        token, created = Token.objects.get_or_create(user=user)
        return Response({"token": token.key, "message": "登录成功！"})
    else:
        return Response({"message": "用户名或密码错误！"}, 401)


# 登出功能由前端实现即可
def logout(request):
    pass


@api_view(["GET"])
def verify(request, **kwargs):
    method = kwargs.get("method")

    if method == "email":
        email = request.query_params.get("email")
        domain = email[email.find("@") + 1:]
        # 检查邮箱是否在白名单内
        if domain not in settings.EMAIL_WHITELIST:
            return Response({"message": "邮箱不在白名单内！"}, 400)
        # 检查用户是否注册
        if User.objects.filter(username=email):
            return Response({"message": "该用户已注册！"}, 400)

        # 设置验证码并发送验证邮件
        verification = random.randint(100000, 999999)
        cache.set(email, verification, settings.VALIDATION_CODE_EXPIRE_TIME * 60)
        mail.delay(
            subject=f'{settings.SITE_NAME} 注册验证',
            content=f'欢迎注册 {settings.SITE_NAME}，您的验证码是: {verification}\r\n验证码的有效期为 {settings.VALIDATION_CODE_EXPIRE_TIME} 分钟\r\n如果您意外地收到了此邮件，请忽略它',
            receivers=[email]
        )
        return Response({'message': '验证邮件发送成功，请查收验证码'})
    else:
        return Response({}, 502)


class RegisterApi(APIView):
    def post(self, request):
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
        if User.objects.filter(username=email).exists():
            return Response({"message": "该用户已注册！"}, 400)

        User.objects.create_user(username=email, password=password)
        return Response({"message": "注册成功！"}, 201)

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
        users = User.objects.filter(username=email)
        if not users:
            return Response({"message": "该用户不存在"}, 400)
        user = users[0]

        user.set_password(password)
        user.save()
        return Response({"message": "已重置密码"}, 200)


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
    reply_to = serializer.validated_data.get('reply_to')
    shadow_text = to_shadow_text(content)

    # 获取匿名信息，如没有则随机选取一个，并判断有无重复
    anonyname = hole.mapping.get(request.user.pk)
    if not anonyname:
        while True:
            anonyname = random.choice(settings.NAME_LIST)
            if anonyname in hole.mapping.values():
                pass
            else:
                hole.mapping[request.user.pk] = anonyname
                break
    # 创建 floor 并增加 hole 的楼层数
    floor = Floor.objects.create(hole=hole, content=content, shadow_text=shadow_text, anonyname=anonyname, user=request.user, reply_to=reply_to)
    hole.reply = hole.reply + 1
    hole.save()
    return hole if category == 'hole' else floor


class HolesApi(APIView):
    permission_classes = [IsAuthenticated, NotSilentOrAdminCanPost, OnlyAdminCanModify]

    def get(self, request, **kwargs):
        # 获取单个
        hole_id = kwargs.get('hole_id')
        if hole_id:
            hole = get_object_or_404(Hole, pk=hole_id)
            Hole.objects.filter(pk=hole_id).update(view=F('view') + 1)
            serializer = HoleSerializer(hole, context={"user": request.user})
            return Response(serializer.data)

        # 获取多个
        start_time = request.query_params.get('start_time')
        length = int(request.query_params.get('length'))
        tag_name = request.query_params.get('tag')

        if tag_name:
            tag = get_object_or_404(Tag, name=tag_name)
            query_set = tag.hole_set.all()
        else:
            query_set = Hole.objects.all()

        holes = query_set.order_by('-time_updated').filter(time_updated__lt=start_time)[:length]
        serializer = HoleSerializer(holes, many=True, context={"user": request.user})
        return Response(serializer.data)

    def post(self, request):
        serializer = HoleSerializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        tag_names = serializer.validated_data.get('tag_names')
        division_id = serializer.validated_data.get('division_id')
        self.check_object_permissions(request, division_id)
        # 实例化 Hole
        hole = Hole(division_id=division_id)
        hole.save()
        # 创建 tag 并添加至 hole
        for tag_name in tag_names:
            tag, created = Tag.objects.get_or_create(name=tag_name)
            hole.tags.add(tag)
        # 保存 hole
        hole.save()

        serializer = HoleSerializer(add_a_floor(request, hole, category='hole'), context={"user": request.user})
        return Response({'message': '发表成功！', 'data': serializer.data}, 201)

    def put(self, request, **kwargs):
        hole_id = kwargs.get('hole_id')
        hole = get_object_or_404(Hole, pk=hole_id)
        serializer = HoleSerializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        tag_names = serializer.validated_data.get('tag_names')
        view = serializer.validated_data.get('view')

        if tag_names:
            hole.tags.clear()
            for tag_name in tag_names:
                tag, created = Tag.objects.get_or_create(name=tag_name)
                hole.tags.add(tag)
        if view:
            hole.view = view

        hole.save()
        serializer = HoleSerializer(hole, context={"user": request.user})
        return Response(serializer.data)

    def delete(self, request, **kwargs):
        # 主题帖不能删除
        return Response(None, 204)


class FloorsApi(APIView):
    permission_classes = [IsAuthenticated, NotSilentOrAdminCanPost, OwnerOrAdminCanModify]

    def get(self, request, **kwargs):
        # 获取单个
        floor_id = kwargs.get('floor_id')
        if floor_id:
            floor = get_object_or_404(Floor, pk=floor_id)
            serializer = FloorSerializer(floor, context={"user": request.user})
            return Response(serializer.data)
        # 获取多个（给定 hole下）
        hole_id = int(request.query_params.get('hole_id'))
        search = request.query_params.get('s')
        query_set = Floor.objects.filter(hole_id=hole_id)
        if search:
            query_set = query_set.filter(shadow_text__icontains=search).order_by('-pk')
        else:
            start_floor = request.query_params.get('start_floor')
            start_floor = int(start_floor) if start_floor else 0
            length = int(request.query_params.get('length'))
            if length:
                query_set = query_set[start_floor: start_floor + length]
            else:
                query_set = query_set[start_floor:]
        serializer = FloorSerializer(query_set, many=True, context={"user": request.user})
        return Response(serializer.data)

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
        floor = get_object_or_404(Floor, pk=floor_id)
        self.check_object_permissions(request, floor)
        if content and content.strip():
            floor.history.append({
                'content': floor.content,
                'altered_by': request.user.pk,
                'altered_time': datetime.now(timezone.utc).isoformat()
            })
            floor.content = content
            floor.shadow_text = to_shadow_text(content)
        if like:
            floor.like_data.append(request.user.pk)
            floor.like += 1
        if fold:
            floor.fold = fold

        floor.save()
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
            floor.shadow_text = '该内容已被作者删除'
        else:  # 管理员删除
            floor.content = delete_reason if delete_reason else '该内容因违反社区规范被删除'
            floor.shadow_text = to_shadow_text(delete_reason) if delete_reason else '该内容因违反社区规范被删除'
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
        query_set = request.user.profile.favorites.all()
        serializer = HoleSerializer(query_set, many=True, context={"user": request.user})
        return Response(serializer.data)

    def post(self, request):
        hole_id = request.data.get('hole_id')
        hole = get_object_or_404(Hole, pk=hole_id)
        profile = get_object_or_404(Profile, user=request.user)
        profile.favorites.add(hole)
        profile.save()
        return Response({'message': '收藏成功'}, 201)

    def put(self, request):
        hole_ids = request.data.get('hole_ids')
        holes = Hole.objects.filter(pk__in=hole_ids)
        profile = get_object_or_404(Profile, user=request.user)
        profile.favorites.set(holes)
        profile.save()
        return Response({'message': '修改成功'}, 200)

    def delete(self, request):
        hole_id = request.data.get('hole_id')
        hole = get_object_or_404(Hole, pk=hole_id)
        profile = get_object_or_404(Profile, user=request.user)
        profile.favorites.remove(hole)
        profile.save()
        return Response({'message': '删除成功'}, 204)


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
            floor.shadow_text = to_shadow_text(delete_reason)
            floor.deleted = True
        if deal.get('silent'):
            profile = floor.user.profile
            current_time_str = profile.permission['silent'].get(str(floor.hole.division_id), '1970-01-01T00:00:00+00:00')
            current_time = parse_datetime(current_time_str)
            expected_time = datetime.now(timezone.utc) + timedelta(days=deal.get('silent'))
            profile.permission['silent'][str(floor.hole.division_id)] = max(current_time, expected_time).isoformat()
            profile.save()

        floor.save()
        report.dealed_by = request.user
        report.dealed = True
        report.save()
        return Response({'message': '举报处理成功'}, 204)

    def put(self, request):
        pass


class ImagesApi(APIView):
    permission_classes = [IsAuthenticated]

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
    permission_classes = [IsAuthenticated]

    def post(self, request):
        from_id = request.data.get('from')
        from_user = get_object_or_404(User, pk=from_id)
        floor = get_object_or_404(Floor, pk=request.data.get('to'))
        to_user = floor.user
        to_id = to_user.pk

        if request.data.get('share_email'):
            message = f'用户看到了你发布的帖子\n{str(floor)}\n希望与你取得联系，TA的邮箱为：{from_user.username}'
            send_message.delay(from_id=from_id, to_id=to_id, message=message)

        elif request.data.get('message'):
            message = request.data.get('message').strip()
            if not is_permitted(request.user, 'admin'):
                return Response(None, 403)
            if not message:
                return Response({'message': 'message不能为空'}, 400)
            send_message.delay(from_id=from_id, to_id=to_id, message=message)

        else:
            return Response(None, 400)

        return Response({'message': '已发送通知'}, 201)
