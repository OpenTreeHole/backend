import os
from random import randint

from django.conf import settings
from django.contrib.auth.hashers import check_password
from django.contrib.auth.models import User
from django.contrib.auth.password_validation import validate_password
from django.core.cache import cache
from django.core.exceptions import ValidationError
from django.shortcuts import get_object_or_404
from rest_framework.authtoken.models import Token
from rest_framework.decorators import api_view
from rest_framework.response import Response

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
            verification = randint(100000, 999999)
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
