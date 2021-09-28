"""
项目初始化时自动创建管理员账户
"""
import os

import django

os.environ.setdefault("DJANGO_SETTINGS_MODULE", "OpenTreeHole.settings")
django.setup()

from django.contrib.auth.models import User

if User.objects.count() == 0:
    print('''
项目初始化时自动创建管理员账户
邮箱为 admin@opentreehole.org，密码为 admin
请尽快登录至管理后台修改管理员信息
''')
    admin = User.objects.create_user('admin@opentreehole.org', password='admin')
    admin.profile.permission['admin'] = '9999-01-01T00:00:00+00:00'
    admin.profile.nickname = 'admin'
    admin.profile.save()
