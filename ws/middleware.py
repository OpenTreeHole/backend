from channels.db import database_sync_to_async
from channels.middleware import BaseMiddleware
from django.contrib.auth.models import AnonymousUser
from rest_framework.authtoken.models import Token


@database_sync_to_async
def get_user(token_key):
    try:
        token = Token.objects.get(key=token_key)
        return token.user
    except Token.DoesNotExist:
        return AnonymousUser()


def find_token_in_headers(headers):
    # scope['headers'] 为二进制编码的元组列表
    for header in headers:
        if header[0].decode() in ('authorization', 'sec-websocket-protocol'):
            return header[1].decode().split(' ')[-1]


def find_in_query_string(query_string, name='token'):
    params = query_string.decode().split('&')
    for param in params:
        if param.startswith(f'{name}='):
            return param.replace(f'{name}=', '')


class TokenAuthMiddleware(BaseMiddleware):
    def __init__(self, inner):
        super().__init__(inner)

    async def __call__(self, scope, receive, send):
        try:
            token_key = find_token_in_headers(scope['headers']) or find_in_query_string(scope['query_string'])
        except ValueError:
            token_key = None
        scope['user'] = AnonymousUser() if token_key is None else await get_user(token_key)
        return await super().__call__(scope, receive, send)
