# Create your tasks here
from __future__ import absolute_import, unicode_literals

import re
import time
from datetime import datetime, timedelta
from smtplib import SMTPException

from celery import shared_task
from celery.schedules import crontab
from django.conf import settings
from django.contrib.auth import get_user_model
from django.core.cache import cache
from django.core.mail import send_mail
from django.db import transaction
from django.db.models import F

from OpenTreeHole.celery import app
from api.models import Hole, ActiveUser
from ws.utils import send_websocket_message_to_group

User = get_user_model()


@shared_task
def hello_world():
    time.sleep(1)
    print('hello world')


@shared_task
def send_email(subject: str, content: str, receivers: list[str], uuid=None) -> dict:
    try:
        send_mail(
            subject=subject,
            message=content,
            from_email=None,
            recipient_list=receivers,
            fail_silently=False,
        )
    except SMTPException as e:
        result = {
            'message': f'邮件发送错误，收件人：{receivers}，错误信息：{e}',
            'success': False
        }
    else:
        result = {
            'message': '邮件发送成功',
            'success': True
        }
    if uuid:
        try:
            send_websocket_message_to_group(uuid, result)
        except:
            print(uuid, type(uuid))
    return result


update_hole_views_pattern = re.compile(r'hole_viewed_(\d+)')
update_last_login_pattern = re.compile(r'user_last_login_(\d+)')


@app.task
def update_hole_views():
    cnt = 0
    with transaction.atomic():
        for key in cache.iter_keys('hole_viewed_*'):
            pattern = update_hole_views_pattern.findall(key)
            if not pattern:
                continue
            cnt += 1
            Hole.objects.filter(pk=int(pattern[0])).update(view=F('view') + cache.get(key, 0))
            cache.delete(key)
    return f'updated {cnt} hole views'


@app.task
def update_last_login():
    cnt = 0
    with transaction.atomic():
        for key in cache.iter_keys('user_last_login_*'):
            pattern = update_last_login_pattern.findall(key)
            if not pattern:
                continue
            cnt += 1
            User.objects.filter(pk=int(pattern[0])).update(last_login=cache.get(key, ''))
            cache.delete(key)
    return f'updated {cnt} user last logins'


@app.task
def calculate_active_user():
    now = datetime.now(settings.TIMEZONE)
    one_day_ago = now - timedelta(days=1)
    one_month_ago = now - timedelta(days=30)
    dau = User.objects.filter(last_login__gt=one_day_ago).count()
    mau = User.objects.filter(last_login__gt=one_month_ago).count()
    obj, created = ActiveUser.objects.update_or_create(
        date=one_day_ago,
        defaults={'dau': dau, 'mau': mau}
    )
    return obj.date, obj.dau


@app.on_after_finalize.connect
def setup_periodic_tasks(sender, **kwargs):
    sender.add_periodic_task(60, update_hole_views.s())  # 每分钟更新一次浏览量
    sender.add_periodic_task(3600, update_last_login.s())  # 每小时更新一次 last_login
    sender.add_periodic_task(crontab(minute=0, hour=0), calculate_active_user.s())  # 每天零点更新日活月活用户数
