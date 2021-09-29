from pathlib import Path

from .config import *

print(f'{SITE_NAME} 正在以开发模式运行，请不要用在生产环境')

# 就是外层的 OpenTreeHole
BASE_DIR = Path(__file__).resolve().parent.parent
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
