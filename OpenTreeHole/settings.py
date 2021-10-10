import os
import sys

if os.environ.get("HOLE_ENV") == "development":
    from OpenTreeHole.development import *
elif os.environ.get("HOLE_ENV") == "production":
    from OpenTreeHole.production import *
else:
    print("未配置 HOLE_ENV 环境变量！请将其配置为 development 或 production")
    from OpenTreeHole.development import *

# Application definition

INSTALLED_APPS = [
    "django.contrib.auth",
    "django.contrib.contenttypes",
    # "django.contrib.sessions",
    # "django.contrib.messages",
    "django.contrib.staticfiles",
    "rest_framework",
    "rest_framework.authtoken",
    'django_celery_results',
    'channels',
    "api.apps.ApiConfig",
]

MIDDLEWARE = [
    # "django.middleware.security.SecurityMiddleware",
    # "django.contrib.sessions.middleware.SessionMiddleware",
    # "django.middleware.common.CommonMiddleware",
    "django.middleware.csrf.CsrfViewMiddleware",
    # "django.contrib.auth.middleware.AuthenticationMiddleware",
    # "django.contrib.messages.middleware.MessageMiddleware",
    # "django.middleware.clickjacking.XFrameOptionsMiddleware",
]

# 开发环境使用内置模板支持一些开发工具
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

ROOT_URLCONF = "api.urls"

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

REST_FRAMEWORK = {
    "DEFAULT_AUTHENTICATION_CLASSES": [
        "rest_framework.authentication.TokenAuthentication"
    ],
    'DEFAULT_RENDERER_CLASSES': [
        'rest_framework.renderers.JSONRenderer'
    ],
    "TEST_REQUEST_DEFAULT_FORMAT": "json",
    'EXCEPTION_HANDLER': 'api.utils.custom_exception_handler',
    'DEFAULT_THROTTLE_CLASSES': [
        'api.throttles.BurstRateThrottle',
        'api.throttles.SustainedRateThrottle',
        'rest_framework.throttling.ScopedRateThrottle',
    ],
    'DEFAULT_THROTTLE_RATES': {
        'burst': THROTTLE_BURST,
        'sustained': THROTTLE_SUSTAINED,
        'email': THROTTLE_EMAIL,
        'upload': THROTTLE_UPLOAD
    }
}
# 测试时不限制访问速率
TESTING = len(sys.argv) > 1 and sys.argv[1] == 'test'
if TESTING:
    del REST_FRAMEWORK['DEFAULT_THROTTLE_RATES']
    del REST_FRAMEWORK['DEFAULT_THROTTLE_CLASSES']

FIXTURE_DIRS = [os.path.join(Path(__file__).resolve().parent, 'fixtures')]

CELERY_RESULT_BACKEND = REDIS_URL
CELERY_BROKER_URL = REDIS_URL

# 自定义用户模型
AUTH_USER_MODEL = 'api.User'

# 遥远的时间
VERY_LONG_TIME = '9999-01-01T00:00:00+00:00'
