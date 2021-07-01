from pathlib import Path

# 就是外层的 OpenTreeHole
BASE_DIR = Path(__file__).resolve().parent.parent
DEBUG = False

# 此处填写你的域名
ALLOWED_HOSTS = ["localhost"]

DATABASES = {
    "default": {
        "ENGINE": "django.db.backends.sqlite3",
        "NAME": BASE_DIR / "db.sqlite3",
    }
}
