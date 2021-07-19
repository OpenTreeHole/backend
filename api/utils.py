from smtplib import SMTPException

from django.core.mail import send_mail
from rest_framework.views import exception_handler


def custom_exception_handler(exc, context):
    # Call REST framework's default exception handler first,
    # to get the standard error response.
    response = exception_handler(exc, context)

    # Now add the HTTP status code to the response.
    if response is not None and response.data.get('detail'):
        response.data['message'] = str(response.data['detail'])
        del (response.data['detail'])

    return response


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
