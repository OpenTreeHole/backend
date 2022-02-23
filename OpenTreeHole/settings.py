import os
import sys

if os.environ.get("HOLE_ENV") == "development":
    from OpenTreeHole.development import *

    HOLE_ENV = "development"
    print(f'{SITE_NAME} 正在以开发模式运行，请不要用在生产环境')

elif os.environ.get("HOLE_ENV") == "testing":
    from OpenTreeHole.testing import *

    HOLE_ENV = "testing"
    print(f'{SITE_NAME} 正在以测试模式运行，请不要用在生产环境')

elif os.environ.get("HOLE_ENV") == "production":
    from OpenTreeHole.production import *

    HOLE_ENV = 'production'
else:
    print("未配置 HOLE_ENV 环境变量！请将其配置为 development / testing / production")
    from OpenTreeHole.development import *

    HOLE_ENV = 'development'
    print(f'{SITE_NAME} 正在以开发模式运行，请不要用在生产环境')

TEMPLATES = [
    {
        "BACKEND": "django.template.backends.django.DjangoTemplates",
        "DIRS": [BASE_DIR / "templates"],
        "APP_DIRS": True,
        "OPTIONS": {
            "context_processors": [
                "django.template.context_processors.debug",
                "django.template.context_processors.request",
                "django.contrib.auth.context_processors.auth",
                "django.contrib.messages.context_processors.messages",
            ],
        },
    },
]

ROOT_URLCONF = "OpenTreeHole.urls"

WSGI_APPLICATION = "OpenTreeHole.wsgi.application"
ASGI_APPLICATION = "OpenTreeHole.asgi.application"

# Password validation
# https://docs.djangoproject.com/en/3.2/ref/settings/#auth-password-validators

AUTH_PASSWORD_VALIDATORS = [
    {
        "NAME": "django.contrib.auth.password_validation.UserAttributeSimilarityValidator",
    },
    {
        "NAME": "django.contrib.auth.password_validation.MinimumLengthValidator",
        "OPTIONS": {
            "min_length": MIN_PASSWORD_LENGTH,
        },
    },
    {
        "NAME": "django.contrib.auth.password_validation.CommonPasswordValidator",
    },
    {
        "NAME": "django.contrib.auth.password_validation.NumericPasswordValidator",
    },
]

# Internationalization
# https://docs.djangoproject.com/en/3.2/topics/i18n/

LANGUAGE_CODE = LANGUAGE

TIME_ZONE = TZ

USE_I18N = True

USE_L10N = True

USE_TZ = True

# Static files (CSS, JavaScript, Images)
# https://docs.djangoproject.com/en/3.2/howto/static-files/

STATIC_URL = "/static/"

# Default primary key field type
# https://docs.djangoproject.com/en/3.2/ref/settings/#default-auto-field

DEFAULT_AUTO_FIELD = "django.db.models.BigAutoField"
DATA_UPLOAD_MAX_MEMORY_SIZE = 30 * 1024 * 1024  # 30M

# 自定义用户模型
AUTH_USER_MODEL = 'api.User'

# 遥远的时间
VERY_LONG_TIME = '9999-01-01T00:00:00+00:00'

CELERY_TIMEZONE = TZ
CELERY_TASK_TIME_LIMIT = 20
CELERY_TASK_COMPRESSION = 'gzip'
CELERY_RESULT_COMPRESSION = 'gzip'
CELERY_WORKER_MAX_TASKS_PER_CHILD = 100
