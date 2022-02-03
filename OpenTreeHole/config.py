"""
站点配置
"""
import json
import os
import uuid
from json import JSONDecodeError
from pathlib import Path

import pytz
from pytz import UnknownTimeZoneError


def get_int_from_env(name: str, value: int) -> int:
    try:
        return int(os.environ.get(name, ''))
    except ValueError:
        return value


def get_array_from_env(name: str, value: list) -> list:
    try:
        return json.loads(os.environ.get(name, ''))
    except JSONDecodeError:
        return value


def get_bool_from_env(name: str, value: bool) -> bool:
    s = os.environ.get(name, '')
    if s == 'True':
        return True
    elif s == 'False':
        return False
    else:
        return value


SITE_NAME = os.environ.get('SITE_NAME', 'Open Tree Hole')  # 网站名称
lower_site_name = SITE_NAME.replace(' ', '').lower()
HOST = os.environ.get('HOST', f'www.{lower_site_name}.com')  # 网站域名
LANGUAGE = os.environ.get('LANGUAGE', 'zh-Hans')  # 语言代码
ALLOW_CONNECT_HOSTS = get_array_from_env('ALLOW_CONNECT_HOSTS', [HOST])  # 允许连接的域名
EMAIL_WHITELIST = get_array_from_env('EMAIL_WHITELIST', ['test.com'])  # 允许注册树洞的邮箱域名
MIN_PASSWORD_LENGTH = get_int_from_env('MIN_PASSWORD_LENGTH', 8)  # 允许的最短用户密码长度
VALIDATION_CODE_EXPIRE_TIME = get_int_from_env('VALIDATION_CODE_EXPIRE_TIME',
                                               5)  # 验证码失效时间（分钟）
MAX_PAGE_SIZE = get_int_from_env('MAX_PAGE_SIZE', 10)
PAGE_SIZE = get_int_from_env('PAGE_SIZE', 10)
FLOOR_PREFETCH_LENGTH = get_int_from_env('FLOOR_PREFETCH_LENGTH', 10)
MAX_TAGS = get_int_from_env('MAX_TAGS', 5)
MAX_TAG_LENGTH = get_int_from_env('MAX_TAG_LENGTH', 16)
TAG_COLORS = get_array_from_env(
    'TAG_COLORS',
    ['red', 'pink', 'purple', 'deep-purple', 'indigo', 'blue', 'light-blue', 'cyan',
     'teal', 'green',
     'light-green', 'lime', 'yellow', 'amber', 'orange', 'deep-orange', 'brown',
     'blue-grey', 'grey']
)
# 时区配置
TZ = os.environ.get('TZ', 'Asia/Shanghai')
try:
    TIMEZONE = pytz.timezone(TZ)
except UnknownTimeZoneError:
    TIMEZONE = pytz.timezone('utc')

# 缓存配置
HOLE_CACHE_SECONDS = get_int_from_env('HOLE_CACHE_SECONDS', 10 * 60)
FLOOR_CACHE_SECONDS = get_int_from_env('FLOOR_CACHE_SECONDS', 10 * 60)

# 访问速率限制
THROTTLE_BURST = os.environ.get('THROTTLE_BURST', '60/min')
THROTTLE_SUSTAINED = os.environ.get('THROTTLE_SUSTAINED', '1000/day')
THROTTLE_EMAIL = os.environ.get('THROTTLE_EMAIL', '30/day')
THROTTLE_UPLOAD = os.environ.get('THROTTLE_UPLOAD', '30/day')

# 数据库配置
DATABASE_HOST = os.environ.get('DATABASE_HOST', 'localhost')
DATABASE_PORT = get_int_from_env('DATABASE_PORT', 3306)
DATABASE_NAME = os.environ.get('DATABASE_NAME', 'hole')
DATABASE_USER = os.environ.get('DATABASE_USER', 'root')
DATABASE_PASSWORD = os.environ.get('DATABASE_PASSWORD', '')
REDIS_URL = os.environ.get('REDIS_URL', 'redis://localhost:6379')

# 邮件配置
EMAIL_HOST = os.environ.get('EMAIL_HOST', '')
EMAIL_PORT = get_int_from_env('EMAIL_PORT', 465)
EMAIL_HOST_USER = os.environ.get('EMAIL_HOST_USER', '')
EMAIL_HOST_PASSWORD = os.environ.get('EMAIL_HOST_PASSWORD', '')
EMAIL_USE_TLS = get_bool_from_env('EMAIL_USE_TLS', False)
EMAIL_USE_SSL = get_bool_from_env('EMAIL_USE_SSL', True)
DEFAULT_FROM_EMAIL = os.environ.get('DEFAULT_FROM_EMAIL', '')  # 默认发件人地址

# 代理配置
HTTP_PROXY = os.environ.get('HTTP_PROXY', None)
# 图片配置
MAX_IMAGE_SIZE = get_int_from_env('MAX_IMAGE_SIZE', 20)  # 最大上传图片大小（MB）
IMAGE_BACKEND = os.environ.get('IMAGE_BACKEND', '')
# Github 图床，具体可参考 https://gitnoteapp.com/zh/extensions/github.html
GITHUB_OWENER = os.environ.get('GITHUB_OWENER', 'OpenTreeHole')
GITHUB_TOKEN = os.environ.get('GITHUB_TOKEN', '123456')
GITHUB_REPO = os.environ.get('GITHUB_REPO', 'images')
GITHUB_BRANCH = os.environ.get('GITHUB_BRANCH', 'master')
# chevereto 图床
CHEVERETO_URL = os.environ.get('CHEVERETO_URL',
                               '')  # e.g. https://www.chevereto.com/api/1/upload
CHEVERETO_TOKEN = os.environ.get('CHEVERETO_TOKEN', '')

# 足够长的密码，供 Django 安全机制
SECRET_KEY = os.environ.get('SECRET_KEY', str(uuid.uuid4()))

# 注册用 API Key (Seed)
REGISTER_API_KEY_SEED = os.environ.get('REGISTER_API_KEY_SEED', 'abcdefg')

# 用户名加密公钥文件(PEM)路径
USERNAME_PUBLIC_KEY_PATH = os.environ.get('USERNAME_PUBLIC_KEY_PATH',
                                          'treehole_demo_public.pem')

USE_REDIS_IN_DEV = get_bool_from_env('USE_REDIS_IN_DEV', False)  # 开发环境中使用 redis

# 推送通知
# Leave APNS_KEY_PATH empty to disable APNS
# NOTE: The APNS KEY must contain both the certificate AND the private key, in PEM format
APNS_KEY_PATH = os.environ.get('APNS_KEY_PATH', '')
APNS_USE_ALTERNATIVE_PORT = get_bool_from_env('APNS_USE_ALTERNATIVE_PORT', False)
MIPUSH_APP_SECRET = os.environ.get(
    'MIPUSH_APP_SECRET', '')  # Leave blank to disable MiPush
PUSH_NOTIFICATION_CLIENT_PACKAGE_NAME_IOS = os.environ.get(
    'PUSH_NOTIFICATION_CLIENT_PACKAGE_NAME_IOS', 'org.opentreehole.client')
PUSH_NOTIFICATION_CLIENT_PACKAGE_NAME_ANDROID = os.environ.get(
    'PUSH_NOTIFICATION_CLIENT_PACKAGE_NAME_ANDROID', 'org.opentreehole.client')

# 就是外层的 OpenTreeHole
BASE_DIR = Path(__file__).resolve().parent.parent

INSTALLED_APPS = [
    "django.contrib.auth",
    "django.contrib.contenttypes",
    "django.contrib.staticfiles",
    "rest_framework",
    "rest_framework.authtoken",
    'django_celery_results',
    'channels',
    "api.apps.ApiConfig",
]

MIDDLEWARE = []

REST_FRAMEWORK = {
    "DEFAULT_AUTHENTICATION_CLASSES": [
        "utils.auth.MyTokenAuthentication"
    ],
    'DEFAULT_RENDERER_CLASSES': [
        'rest_framework.renderers.JSONRenderer'
    ],
    "TEST_REQUEST_DEFAULT_FORMAT": "json",
    'EXCEPTION_HANDLER': 'utils.exception.custom_exception_handler',
}

if __name__ == '__main__':
    print(get_array_from_env('test', [1]))
