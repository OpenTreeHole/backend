import re

from asgiref.sync import async_to_sync
from channels.layers import get_channel_layer
from rest_framework.views import exception_handler


def send_message_to_user(user, content):
    """
    向用户发送消息
    Args:
        user: 用户对象
        content: 消息内容

    Returns: None

    """
    channel_layer = get_channel_layer()
    async_to_sync(channel_layer.group_send)(
        f'user-{user.id}',  # Channels 组名称
        {
            "type": "notification",
            "content": content,
        }
    )


def custom_exception_handler(exc, context):
    # Call REST framework's default exception handler first,
    # to get the standard error response.
    response = exception_handler(exc, context)

    # 默认错误消息字段改为“message”
    if response is not None and response.data.get('detail'):
        response.data['message'] = str(response.data['detail'])
        del (response.data['detail'])

    return response


def to_shadow_text(content):
    return re.sub(r'([\s#*_!>`$|:,\-\[\]-]|\d+\.|\(.+?\)|<.+?>)', '', content)
