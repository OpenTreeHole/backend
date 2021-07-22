from pathlib import Path

from .config import *

# 就是外层的 OpenTreeHole
BASE_DIR = Path(__file__).resolve().parent.parent
DEBUG = False

# 此处填写你的域名
ALLOWED_HOSTS = ALLOW_CONNECT_HOSTS

# 生产环境使用 Mysql 数据库
DATABASES = {
    "default": {
        "ENGINE": "django.db.backends.mysql",
        "NAME": DATABASE_NAME,
        "USER": DATABASE_USER,
        "PASSWORD": DATABASE_PASSWORD,
        "HOST": DATABASE_HOST,
        "PORT": DATABASE_PORT,
    }
}

# 生产环境使用 Redis 作为缓存
CACHES = {
    "default": {
        "BACKEND": "django_redis.cache.RedisCache",
        "LOCATION": "redis://127.0.0.1:6379",
        "OPTIONS": {
            "CLIENT_CLASS": "django_redis.client.DefaultClient",
        },
    },
}

# 生产环境SMTP发送邮件
EMAIL_BACKEND = 'django.core.mail.backends.smtp.EmailBackend'
