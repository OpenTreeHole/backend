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
    def _authenticate(self, auth_method, token):
        if auth_method == 'token':
            return self.authenticate_credentials(token)
        elif auth_method == 'bearer':
            try:
                payload = jwt.decode(token, verify=False, options={"verify_signature": False})
            except jwt.DecodeError:
                raise AuthenticationFailed('jwt token invalid')
            _id = payload.get('uid')
            try:
                user = User.objects.get(id=_id)
            except User.DoesNotExist:
                email = f'user#{_id}@fduhole.com'
                user = User.objects.create_user(id=_id, email=email, password=email)
            return user, token

    def authenticate(self, request):
        if request.headers.get('x-anonymous-consumer'):
            return
        auth = get_authorization_header(request).split()
        if len(auth) != 2:
            return
        try:
            uid = int(request.headers.get('x-consumer-username'))
        except:
            return
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
        return user, auth[1].decode()


async def async_token_auth(request):
    method = MyTokenAuthentication().authenticate
    try:
        user, token = await sync_to_async(method)(request)
    except (AuthenticationFailed, TypeError):
        request.user = AnonymousUser()
        return request
    request.user = user
    return request
