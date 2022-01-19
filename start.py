"""
项目初始化时自动创建管理员账户
"""
import os

import django
from django.contrib.auth import get_user_model

os.environ.setdefault("DJANGO_SETTINGS_MODULE", "OpenTreeHole.settings")
django.setup()

from api.models import Division

User = get_user_model()
if User.objects.count() == 0:
    print('''
项目初始化时自动创建管理员账户
邮箱为 admin@opentreehole.org，密码为 admin
请尽快登录至管理后台修改管理员信息
''')
    User.objects.create_superuser(email='admin@opentreehole.org', password='admin')

if Division.objects.count() == 0:
    Division.objects.create(name='树洞', description='论坛板块')
    Division.objects.create(name='表白墙', description='道出心声')
    Division.objects.create(name='评教', description='分享一些经验')
    Division.objects.create(name='站务', description='论坛管理')
