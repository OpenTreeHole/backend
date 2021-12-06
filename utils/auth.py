"""
注册、登录、用户管理
"""
import base64
import hashlib

from Crypto.Cipher import PKCS1_v1_5 as PKCS1_cipher
from Crypto.PublicKey import RSA
from django.conf import settings


def check_api_key(key_to_check):
    return key_to_check == settings.REGISTER_API_KEY_SEED


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


def sha512(string):
    return hashlib.sha512(bytes(string.encode('utf-8'))).hexdigest()


if __name__ == '__main__':
    encrypted = encrypt_email('hi')
    print(decrypt_email(encrypted))
