from smtplib import SMTPException

from django.core.mail import send_mail


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
        return {'message': '邮件发送错误，收件人：{}，错误信息：{}'.format(receivers, e), 'code': 502}
    else:
        return {'message': '邮件发送成功！', 'code': 200}
