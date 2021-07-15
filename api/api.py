import os
import random

from django.contrib.auth.hashers import check_password
from django.contrib.auth.password_validation import validate_password
from django.core.cache import cache
from django.core.exceptions import ValidationError
from django.shortcuts import get_object_or_404
from django.utils.dateparse import parse_datetime
from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.views import APIView

from api.models import *
from api.serializers import *
from api.utils import mail


# 发送 csrf 令牌
# from django.views.decorators.csrf import ensure_csrf_cookie
# @method_decorator(ensure_csrf_cookie)


@api_view(["GET"])
def index(request):
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
        elif User.objects.filter(username=email):
            return Response({"message": "该用户已注册！"}, 400)
        # 正确无误，发送验证邮件
        else:
            verification = random.randint(100000, 999999)
            cache.set(email, verification, settings.VALIDATION_CODE_EXPIRE_TIME * 60)
            # 开发环境不发送邮件
            if os.environ.get('ENV') == 'development':
                return Response({})
            elif os.environ.get('ENV') == 'production':
                mail_result = mail(
                    subject='{} 注册验证'.format(settings.SITE_NAME),
                    content='欢迎注册 {}，您的验证码是: {}\r\n验证码的有效期为 {} 分钟\r\n如果您意外地收到了此邮件，请忽略它'
                        .format(settings.SITE_NAME, verification, settings.VALIDATION_CODE_EXPIRE_TIME),
                    receivers=[email]
                )
                return Response({'message': mail_result['message']}, mail_result['code'])
            else:
                return Response({}, 502)


@api_view(["POST"])
def register(request):
    email = request.data.get("email")
    password = request.data.get("password")
    verification = request.data.get("verification")

    if not verification:
        return Response({"message": "验证码不能为空！"}, 400)
    # 转义表单数据（应该使用 JSON）
    try:
        verification = int(verification)
    except TypeError:
        return Response({"message": "验证码格式错误！"}, 400)
    try:
        validate_password(password)
    except ValidationError as e:
        return Response({'message': '\n'.join(e)}, 400)

    if not cache.get(email) or not verification or not cache.get(email) == verification:
        return Response({"message": "注册校验未通过！"}, 400)
    else:
        User.objects.create_user(username=email, password=password)
        return Response({"message": "注册成功！"})


def add_a_floor(request, hole, type):
    """
    增加一条回复帖
    Args:
        request:
        hole:       hole对象

    Returns:        Response

    """
    # 校验 content
    content = request.data.get('content')
    if not content:
        return Response({'message': '内容不能为空'}, 400)
    content = content.strip()
    if not content:
        return Response({'message': '内容不能为空'}, 400)
    # 校验 reply_to
    reply_to = request.data.get('reply_to')

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
    floor = Floor.objects.create(hole=hole, content=content, anonyname=anonyname, user=request.user, reply_to=reply_to)
    hole.reply = hole.reply + 1
    hole.save()
    return hole if type == 'hole' else floor


class HolesApi(APIView):
    def post(self, request):
        serializer = HoleSerializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        tag_names = serializer.validated_data.get('tag_names')
        division_id = serializer.validated_data.get('division_id')

        # 实例化 Hole
        hole = Hole(division_id=division_id)
        hole.save()
        # 创建 tag 并添加至 hole
        for tag_name in tag_names:
            tag, created = Tag.objects.get_or_create(name=tag_name)
            hole.tags.add(tag)
        # 保存 hole
        hole.save()

        serializer = HoleSerializer(add_a_floor(request, hole, type='hole'), context={"user": request.user})
        return Response({'message': '发表成功！', 'data': serializer.data}, 201)

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


class FloorsApi(APIView):
    def post(self, request):
        hole_id = request.data.get('hole_id')
        hole = get_object_or_404(Hole, pk=hole_id)
        serializer = FloorSerializer(add_a_floor(request, hole, type='floor'), context={"user": request.user})
        return Response({'message': '发表成功！', 'data': serializer.data}, 201)
