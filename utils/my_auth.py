from datetime import datetime

import jwt
from asgiref.sync import sync_to_async
from django.conf import settings
from django.contrib.auth.models import AnonymousUser
from django.core.cache import cache
from rest_framework.authentication import TokenAuthentication, get_authorization_header
from rest_framework.exceptions import AuthenticationFailed

from api.models import User


class MyTokenAuthentication(TokenAuthentication):
    def authenticate(self, request):
        if request.headers.get('x-anonymous-consumer'):
            return
        auth = get_authorization_header(request).split()
        if len(auth) != 2:
            return
        token = auth[1].decode()
        uid = request.headers.get('x-consumer-username')
        if not uid:
            user, token = self.authenticate_credentials(token)
        else:
            try:
                user = User.objects.get(id=uid)
            except User.DoesNotExist:
                email = f'user#{uid}@fduhole.com'
                user = User.objects.create_user(id=uid, email=email, password=email)
        cache.set(
            f'user_last_login_{user.id}',
            datetime.now(settings.TIMEZONE).isoformat(),
            86400
        )
        return user, token


async def async_token_auth(request):
    method = MyTokenAuthentication().authenticate
    try:
        user, token = await sync_to_async(method)(request)
    except (AuthenticationFailed, TypeError):
        request.user = AnonymousUser()
        return request
    request.user = user
    return request
