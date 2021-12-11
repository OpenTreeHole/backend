"""
站点配置
"""

import os
import uuid
from pathlib import Path

SITE_NAME = 'Open Tree Hole'  # 网站名称
lower_site_name = SITE_NAME.replace(' ', '').lower()
HOST = f'www.{lower_site_name}.com'  # 网站域名
TZ = "Asia/Shanghai"  # 时区
LANGUAGE = "zh-Hans"  # 语言代码
ALLOW_CONNECT_HOSTS = [HOST]  # 允许连接的域名
EMAIL_WHITELIST = ["test.com"]  # 允许注册树洞的邮箱域名
MIN_PASSWORD_LENGTH = 8  # 允许的最短用户密码长度
VALIDATION_CODE_EXPIRE_TIME = 5  # 验证码失效时间（分钟）
MAX_PAGE_SIZE = 10
PAGE_SIZE = 10
FLOOR_PREFETCH_LENGTH = 10
MAX_TAGS = 5
MAX_TAG_LENGTH = 8
TAG_COLORS = ('red', 'pink', 'purple', 'deep-purple', 'indigo', 'blue', 'light-blue', 'cyan', 'teal', 'green',
              'light-green', 'lime', 'yellow', 'amber', 'orange', 'deep-orange', 'brown', 'blue-grey', 'grey')

# 缓存配置
HOLE_CACHE_SECONDS = 10 * 60
FLOOR_CACHE_SECONDS = 10 * 60

# 访问速率限制
THROTTLE_BURST = '10/min'
THROTTLE_SUSTAINED = '1000/day'
THROTTLE_EMAIL = '30/day'
THROTTLE_UPLOAD = '30/day'

# 数据库配置
DATABASE_HOST = "localhost"  # 数据库主机
DATABASE_PORT = 3306  # 数据库端口
DATABASE_NAME = "hole"  # 数据库名称
DATABASE_USER = "root"  # 数据库用户
DATABASE_PASSWORD = ""  # 数据库密码
REDIS_URL = 'redis://localhost:6379'  # redis 缓存地址

# 邮件配置
EMAIL_HOST = ''
EMAIL_PORT = 465
EMAIL_HOST_USER = ''
EMAIL_HOST_PASSWORD = ''
EMAIL_USE_TLS = False
EMAIL_USE_SSL = True
DEFAULT_FROM_EMAIL = ''  # 默认发件人地址

# 图片配置
MAX_IMAGE_SIZE = 20  # 最大上传图片大小（MB）

# 采用 Github 图床，具体可参考 https://gitnoteapp.com/zh/extensions/github.html
GITHUB_OWENER = 'OpenTreeHole'
GITHUB_TOKEN = '123456'
GITHUB_REPO = 'images'
GITHUB_BRANCH = 'master'

# 足够长的密码，供 Django 安全机制
SECRET_KEY = str(uuid.uuid4())

# 注册用 API Key (Seed)
REGISTER_API_KEY_SEED = "abcdefg"

# 用户名加密公钥文件(PEM)路径
USERNAME_PUBLIC_KEY_PATH = "treehole_demo_public.pem"

USE_REDIS_IN_DEV = False  # 开发环境中使用 redis

# 推送通知
# Leave APNS_KEY_PATH empty to disable APNS
APNS_KEY_PATH = ""  # NOTE: The APNS KEY must contain both the certificate AND the private key, in PEM format
APNS_USE_ALTERNATIVE_PORT = False
MIPUSH_APP_SECRET = ""  # Leave blank to disable MiPush
PUSH_NOTIFICATION_CLIENT_PACKAGE_NAME_IOS = "org.opentreehole.client"
PUSH_NOTIFICATION_CLIENT_PACKAGE_NAME_ANDROID = "org.opentreehole.client"

# 就是外层的 OpenTreeHole
BASE_DIR = Path(__file__).resolve().parent.parent

# 用环境变量中的配置覆盖
envs = os.environ
local = locals().copy()  # 拷贝一份，否则运行时 locals() 会改变
for item in local:
    if item.startswith('_') or item == 'os':  # 内置变量名不考虑
        continue
    if item in envs:
        try:
            exec(f'{item} = eval(envs.get(item))')  # 非字符串类型使用 eval() 转换
        except:
            exec(f'{item} = envs.get(item)')  # 否则直接为字符串
