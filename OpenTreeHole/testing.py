# noinspection PyUnresolvedReferences
from .config import *

# noinspection PyUnresolvedReferences
from .development import DEBUG, ALLOWED_HOSTS, REST_FRAMEWORK

# noinspection PyUnresolvedReferences
from .production import INSTALLED_APPS, MIDDLEWARE, DATABASES, CACHES, EMAIL_BACKEND, CHANNEL_LAYERS, CELERY_BROKER_URL, CELERY_RESULT_BACKEND
