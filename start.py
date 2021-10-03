"""
项目初始化时自动创建管理员账户
"""
import os

import django
from django.contrib.auth import get_user_model

os.environ.setdefault("DJANGO_SETTINGS_MODULE", "OpenTreeHole.settings")
django.setup()

User = get_user_model()
if User.objects.count() == 0:
    print('''
项目初始化时自动创建管理员账户
邮箱为 admin@opentreehole.org，密码为 admin
请尽快登录至管理后台修改管理员信息
''')
    User.objects.create_superuser(email='admin@opentreehole.org', password='admin')
