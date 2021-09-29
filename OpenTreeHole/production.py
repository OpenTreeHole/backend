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
        "ENGINE": "mysql.connector.django",
        "NAME": DATABASE_NAME,
        "USER": DATABASE_USER,
        "PASSWORD": DATABASE_PASSWORD,
        "HOST": DATABASE_HOST,
        "PORT": DATABASE_PORT,
        'OPTIONS': {
            'auth_plugin': 'mysql_native_password'
        }
    }
}

# 生产环境使用 Redis 作为缓存
CACHES = {
    "default": {
        "BACKEND": "django_redis.cache.RedisCache",
        "LOCATION": REDIS_ADDRESS,
        "OPTIONS": {
            "CLIENT_CLASS": "django_redis.client.DefaultClient",
        },
    },
}

# 生产环境SMTP发送邮件
EMAIL_BACKEND = 'django.core.mail.backends.smtp.EmailBackend'

# channels 通道层，使用 redis
CHANNEL_LAYERS = {
    "default": {
        "BACKEND": "channels_redis.core.RedisChannelLayer",
        "CONFIG": {
            "hosts": [REDIS_ADDRESS],
        },
    },
}
