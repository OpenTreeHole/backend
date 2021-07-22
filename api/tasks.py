# Create your tasks here
from __future__ import absolute_import, unicode_literals

import os
import time
from smtplib import SMTPException

import httpx
from celery import shared_task
from django.core.mail import send_mail


@shared_task
def hello_world():
    time.sleep(1)
    print('hello world')


@shared_task
def mail(subject, content, receivers):
    try:
        send_mail(
            subject=subject,
            message=content,
            from_email=None,
            recipient_list=receivers,
            fail_silently=False,
        )
    except SMTPException as e:
        return '邮件发送错误，收件人：{}，错误信息：{}'.format(receivers, e)
    else:
        return '邮件发送成功！'


@shared_task
def post_image_to_github(url, headers, body):
    proxies = 'http://localhost:7890' if os.environ.get("ENV") == "development" else None
    with httpx.Client(proxies=proxies) as client:
        r = client.put(url, headers=headers, json=body)
    if r.status_code == 201:
        return '上传成功'
    else:
        return '上传失败', r.json()
