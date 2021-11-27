"""
注册、登录、用户管理
"""

from django.conf import settings

PKCS1_PUBLIC_CIPHER = None


def check_api_key(key_to_check):
    return key_to_check == settings.REGISTER_API_KEY_SEED


def encrypt_email(email_cleartext):
    """
    Provide basic encryption

    To decrypt, use:
    print(PKCS1_PUBLIC_CIPHER.decrypt(base64.b64decode(encoded.encode("utf8"))).decode("utf8"))
    """
    # global PKCS1_PUBLIC_CIPHER
    # if PKCS1_PUBLIC_CIPHER is None:
    #     with open(settings.USERNAME_PUBLIC_KEY_PATH, 'r') as file:
    #         PKCS1_PUBLIC_CIPHER = PKCS1_OAEP.new(RSA.importKey(file.read()))
    # return base64.b64encode(PKCS1_PUBLIC_CIPHER.encrypt(email_cleartext.encode("utf8"))).decode("utf8")
    return email_cleartext
