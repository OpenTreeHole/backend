# Create your tasks here
from __future__ import absolute_import, unicode_literals

import os
import time
from smtplib import SMTPException

import httpx
from celery import shared_task
from django.core.cache import cache
from django.core.mail import send_mail
from django.db.models import F

from OpenTreeHole.celery import app
from api.models import Hole
from ws.utils import send_websocket_message_to_group


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


@shared_task
def post_image_to_github(url: str, headers: dict, body: dict, user_id: int) -> dict:
    proxies = os.environ.get("HTTP_PROXY")
    with httpx.Client(proxies=proxies) as client:
        r = client.put(url, headers=headers, json=body)
    if r.status_code == 201:
        result = {'message': '上传成功', 'success': True}
    else:
        result = {'message': '上传失败', 'data': r.json(), 'success': False}
    send_websocket_message_to_group(f'user-{user_id}', result)
    return result


@app.task
def update_hole_views():
    cached = cache.get('hole_views', {})
    for id in cached:
        if cached[id] > 0:
            Hole.objects.filter(pk=id).update(view=F('view') + cached[id])
            cached[id] = 0
    cache.set('hole_views', cached, None)


@app.on_after_finalize.connect
def setup_periodic_tasks(sender, **kwargs):
    sender.add_periodic_task(60, update_hole_views.s())  # 每分钟更新一次浏览量
