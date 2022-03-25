"""
注册、登录、用户管理
"""
import base64
import hashlib
import time
from datetime import datetime

import jwt
import pyotp
from Crypto.Cipher import PKCS1_v1_5 as PKCS1_cipher
from Crypto.PublicKey import RSA
from asgiref.sync import sync_to_async
from django.conf import settings
from django.contrib.auth import get_user_model
from django.contrib.auth.models import AnonymousUser
from django.core.cache import cache
from rest_framework.authentication import TokenAuthentication, get_authorization_header
from rest_framework.exceptions import AuthenticationFailed

User = get_user_model()


async def async_token_auth(request):
    method = MyTokenAuthentication().authenticate
    try:
        user, token = await sync_to_async(method)(request)
    except (AuthenticationFailed, TypeError):
        request.user = AnonymousUser()
        return request
    request.user = user
    return request


class MyTokenAuthentication(TokenAuthentication):
    def _authenticate(self, auth_method, token):
        if auth_method == 'token':
            return self.authenticate_credentials(token)
        elif auth_method == 'bearer':
            try:
                payload = jwt.decode(token, options={"verify_signature": False})
            except jwt.DecodeError:
                raise AuthenticationFailed('jwt token invalid')
            try:
                user = User.objects.get(id=payload.get('uid'))
            except User.DoesNotExist:
                user = User.objects.create(id=payload.get('uid'), email='', password='', identifier='')
            return user, token

    def authenticate(self, request):
        auth = get_authorization_header(request).split()
        if len(auth) != 2:
            raise AuthenticationFailed('token invalid')
        authenticated = self._authenticate(
            auth_method=auth[0].decode().lower(),
            token=auth[1].decode()
        )
        if authenticated:
            user, token = authenticated
            cache.set(
                f'user_last_login_{user.id}',
                datetime.now(settings.TIMEZONE).isoformat(),
                86400
            )
        return authenticated


apikey_verifier_totp = pyotp.TOTP(
    str(base64.b32encode(bytearray(settings.REGISTER_API_KEY_SEED, 'ascii')).decode(
        'utf-8')), digest=hashlib.sha256, interval=5, digits=16)


def check_api_key(key_to_check):
    return apikey_verifier_totp.verify(key_to_check, valid_window=1)


def get_key(key_file):
    """
    generate private key:
        openssl genrsa -out treehole_demo_private.pem 4096

    generate public key:
        openssl rsa -in  treehole_demo_private.pem -pubout -out  treehole_demo_public.pem
    """
    with open(key_file) as f:
        data = f.read()
        key = RSA.importKey(data)
    return key


PUBLIC_KEY = get_key(settings.USERNAME_PUBLIC_KEY_PATH)

CIPHER = PKCS1_cipher.new(PUBLIC_KEY)


def encrypt_email(plaintext):
    """
    RSA encryption
    """
    # To decrypt, use:
    # print(PKCS1_PUBLIC_CIPHER.decrypt(base64.b64decode(encoded.encode("utf8"))).decode("utf8"))

    # global PKCS1_PUBLIC_CIPHER
    # if PKCS1_PUBLIC_CIPHER is None:
    #     with open(settings.USERNAME_PUBLIC_KEY_PATH, 'r') as file:
    #         PKCS1_PUBLIC_CIPHER = PKCS1_OAEP.new(RSA.importKey(file.read()))
    # return base64.b64encode(PKCS1_PUBLIC_CIPHER.encrypt(email_cleartext.encode("utf8"))).decode("utf8")

    encrypted_bytes = CIPHER.encrypt(bytes(plaintext.encode('utf-8')))

    return base64.b64encode(encrypted_bytes).decode('utf-8')


def decrypt_email(encrypted):
    """
    Only for test use
    """
    private_key = get_key('treehole_demo_private.pem')
    cipher = PKCS1_cipher.new(private_key)
    back_text = cipher.decrypt(base64.b64decode(encrypted), 0)
    return back_text.decode('utf-8')


def sha512(string: str) -> str:
    byte_string = bytes(string.encode('utf-8'))
    return hashlib.sha512(byte_string).hexdigest()


def many_hashes(string: str) -> str:
    iterations = 1
    byte_string = bytes(string.encode('utf-8'))
    return hashlib.pbkdf2_hmac('sha3_512', byte_string, b'', iterations).hex()


if __name__ == '__main__':
    start = time.time()

    many_hashes('hi')

    end = time.time()
    print(end - start)
