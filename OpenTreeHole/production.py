from .config import *

DEBUG = False

# 此处填写你的域名
ALLOWED_HOSTS = ALLOW_CONNECT_HOSTS

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

MIDDLEWARE = [
    "django.middleware.csrf.CsrfViewMiddleware",
]

# 生产环境使用 Mysql 数据库
DATABASES = {
    "default": {
        "ENGINE": "django.db.backends.mysql",
        "NAME": DATABASE_NAME,
        "USER": DATABASE_USER,
        "PASSWORD": DATABASE_PASSWORD,
        "HOST": DATABASE_HOST,
        "PORT": DATABASE_PORT,
        'OPTIONS': {
            'auth_plugin': 'mysql_native_password',
            'charset': 'utf8mb4'
        }
    }
}

# 生产环境使用 Redis 作为缓存
CACHES = {
    "default": {
        "BACKEND": "django_redis.cache.RedisCache",
        "LOCATION": REDIS_URL,
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
            "hosts": [REDIS_URL],
        },
    },
}

REST_FRAMEWORK = {
    "DEFAULT_AUTHENTICATION_CLASSES": [
        "rest_framework.authentication.TokenAuthentication"
    ],
    'DEFAULT_RENDERER_CLASSES': [
        'rest_framework.renderers.JSONRenderer'
    ],
    "TEST_REQUEST_DEFAULT_FORMAT": "json",
    'EXCEPTION_HANDLER': 'utils.exception.custom_exception_handler',
    'DEFAULT_THROTTLE_CLASSES': [
        'utils.throttles.BurstRateThrottle',
        'utils.throttles.SustainedRateThrottle',
        'rest_framework.throttling.ScopedRateThrottle',
    ],
    'DEFAULT_THROTTLE_RATES': {
        'burst': THROTTLE_BURST,
        'sustained': THROTTLE_SUSTAINED,
        'email': THROTTLE_EMAIL,
        'upload': THROTTLE_UPLOAD
    }
}

# celery 使用 redis
CELERY_BROKER_URL = REDIS_URL
CELERY_RESULT_BACKEND = REDIS_URL
