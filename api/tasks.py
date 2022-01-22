# Create your tasks here
from __future__ import absolute_import, unicode_literals

import re
import time
from smtplib import SMTPException

from celery import shared_task
from django.contrib.auth import get_user_model
from django.core.cache import cache
from django.core.mail import send_mail
from django.db import transaction
from django.db.models import F

from OpenTreeHole.celery import app
from api.models import Hole
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


@app.task
def update_hole_views():
    cached = cache.get('hole_views', {})
    for id in cached:
        if cached[id] > 0:
            Hole.objects.filter(pk=id).update(view=F('view') + cached[id])
            cached[id] = 0
    cache.set('hole_views', cached, None)


@app.task
def update_last_login():
    with transaction.atomic():
        for key in cache.iter_keys('user_last_login_*'):
            pattern = re.findall(r'user_last_login_(\d+)', key)
            if not pattern:
                continue
            User.objects.filter(pk=int(pattern[0])).update(last_login=cache.get(key, ''))
            cache.delete(key)


@app.on_after_finalize.connect
def setup_periodic_tasks(sender, **kwargs):
    sender.add_periodic_task(60, update_hole_views.s())  # 每分钟更新一次浏览量
    sender.add_periodic_task(3600, update_last_login.s())  # 每小时更新一次 last_login
