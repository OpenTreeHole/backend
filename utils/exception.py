from rest_framework.exceptions import APIException
from rest_framework.views import exception_handler


def custom_exception_handler(exc, context):
    # Call REST framework's default exception handler first,
    # to get the standard error response.
    response = exception_handler(exc, context)

    # 默认错误消息字段改为“message”
    if response is not None and response.data.get('detail'):
        response.data['message'] = str(response.data['detail'])
        del (response.data['detail'])

    return response


class BadRequest(APIException):
    status_code = 400
    default_detail = 'Invalid input.'


class Forbidden(APIException):
    status_code = 403
    default_detail = 'Forbidden'


class ServerError(APIException):
    status_code = 500
    default_detail = 'Internal Server Error'
