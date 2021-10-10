from .config import *

DEBUG = True

# 允许任意 host
ALLOWED_HOSTS = ['*']

# 开发环境使用 sqlite 数据库
DATABASES = {
    "default": {
        "ENGINE": "django.db.backends.sqlite3",
        "NAME": BASE_DIR / "db.sqlite3",
    }
}

# 开发环境使用本地内存作为缓存
CACHES = {
    "default": {
        "BACKEND": "django.core.cache.backends.locmem.LocMemCache",
        "LOCATION": "unique-snowflake",
    }
}

# 开发环境邮件发送至控制台
EMAIL_BACKEND = 'django.core.mail.backends.console.EmailBackend'

# channels 通道层，使用内存
CHANNEL_LAYERS = {
    "default": {
        "BACKEND": "channels.layers.InMemoryChannelLayer"
    }
}

# 开发环境不限制 API 访问速率
REST_FRAMEWORK = {
    "DEFAULT_AUTHENTICATION_CLASSES": [
        "rest_framework.authentication.TokenAuthentication"
    ],
    'DEFAULT_RENDERER_CLASSES': [
        'rest_framework.renderers.JSONRenderer'
    ],
    "TEST_REQUEST_DEFAULT_FORMAT": "json",
    'EXCEPTION_HANDLER': 'api.utils.custom_exception_handler',
}

# LOGGING = {
#     'version': 1,
#     'disable_existing_loggers': False,
#     'handlers': {
#         'console': {
#             'level': 'DEBUG',
#             'class': 'logging.StreamHandler',
#         },
#     },
#     'loggers': {
#         'django.db.backends': {  # 在终端打印 sql 语句
#             'handlers': ['console'],
#             'propagate': True,
#             'level': 'DEBUG',
#         },
#     }
# }
