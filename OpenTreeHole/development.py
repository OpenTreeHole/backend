import sqlite3

from .config import *

DEBUG = True

# 允许任意 host
ALLOWED_HOSTS = ['*']

INSTALLED_APPS = INSTALLED_APPS + ['silk']

MIDDLEWARE = MIDDLEWARE + [
    'django.contrib.sessions.middleware.SessionMiddleware',
    'silk.middleware.SilkyMiddleware',
]

# 开发环境使用 sqlite 数据库
DATABASES = {
    "default": {
        "ENGINE": "django.db.backends.sqlite3",
        "NAME": BASE_DIR / "db.sqlite3",
    }
}

if USE_REDIS_IN_DEV:
    # noinspection PyUnresolvedReferences
    from .production import CACHES, CHANNEL_LAYERS, CELERY_BROKER_URL, \
        CELERY_RESULT_BACKEND
else:
    # 开发环境使用本地内存作为缓存
    CACHES = {
        "default": {
            "BACKEND": "django.core.cache.backends.locmem.LocMemCache",
            "LOCATION": "unique-snowflake",
        }
    }

    # channels 通道层，使用内存
    CHANNEL_LAYERS = {
        "default": {
            "BACKEND": "channels.layers.InMemoryChannelLayer"
        }
    }

    # celery 使用 sqlite
    sqlite3.connect('celery.sqlite3')
    CELERY_BROKER_URL = 'sqla+sqlite:///celery.sqlite3'

# 开发环境邮件发送至控制台
EMAIL_BACKEND = 'django.core.mail.backends.console.EmailBackend'

# 开发环境不限制 API 访问速率
REST_FRAMEWORK = REST_FRAMEWORK

# silk profiling
SILKY_PYTHON_PROFILER = True
