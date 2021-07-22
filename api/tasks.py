# Create your tasks here
from __future__ import absolute_import, unicode_literals

import time
from smtplib import SMTPException

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
