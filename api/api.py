from django.contrib.auth.models import User
from django.contrib.auth.hashers import make_password, check_password
from django.shortcuts import get_object_or_404

from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.views import APIView
from rest_framework.authtoken.models import Token

# 发送 csrf 令牌
# from django.views.decorators.csrf import ensure_csrf_cookie
# @method_decorator(ensure_csrf_cookie)


@api_view(["GET"])
def index(request):
    return Response({"message": "Hello world!"})


@api_view(["POST"])
def login(request):
    username = request.data.get("username")
    password = request.data.get("password")

    user = get_object_or_404(User, username=username)
    if check_password(password, user.password):
        token, created = Token.objects.get_or_create(user=user)
        return Response({"token": token.key})
    else:
        return Response({"message": "用户名或密码错误！"}, 401)


# 登出功能由前端实现即可
def logout(request):
    pass
