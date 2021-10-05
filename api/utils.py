import collections
from os import environ
import re

from asgiref.sync import async_to_sync
from channels.layers import get_channel_layer
from rest_framework.views import exception_handler

from api.models import Message
from api.serializers import MessageSerializer

from apns2.client import APNsClient
from apns2.payload import Payload as APNsPayload
from OpenTreeHole.config import PUSH_NOTIFICATION_CLIENT_PACKAGE_NAME_IOS, PUSH_NOTIFICATION_CLIENT_PACKAGE_NAME_ANDROID

# APNS global definition
Notification = collections.namedtuple('Notification', ['token', 'payload'])
# NOTE: The apns_key.pem must contain both the certificate AND the private key
apns_client = APNsClient('apns_key.pem', use_sandbox=(environ.get("HOLE_ENV") != "production"),
                         use_alternative_port=False)


class MessageSender:
    """
    批量发送消息助手

    使用方法：
    - 首先调用message_sender.create_and_queue_or_send_message(user, message, data, 'favorite')
        此时会自动判断消息是否需要批量发送，如果不需要（例如WebSocket），则直接发送
        如果需要（例如APNS），则加入队列，后期调用commit()发送
    - 最后调用message_sender.commit()批量发送队列中的消息
    """
    apns_notifications = []

    def __init__(self, user=None, message=None, data=None, code='') -> None:
        if user and message:
            if data is None:
                data = {}
            self.create_and_queue_or_send_message(user, message, data, code)

    def __send_websocket_message_to_user(self, user, content):
        """
        向用户发送 WebSocket 消息
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

    def create_and_queue_or_send_message(self, user, message, data=None, code=''):
        """
        首先在数据库中创建Message
        如果消息需要打包发送(apns)，则加入队列
        如果不需要打包发送(websocket)，则直接发送
        """
        if data is None:
            data = {}

        instance = Message.objects.create(user=user, message=message, data=data, code=code)
        payload = MessageSerializer(instance).data

        # WS
        self.__send_websocket_message_to_user(user, payload)

        # APNS
        apns_payload = APNsPayload(alert=message, sound="default", badge=1, thread_id=code, custom=payload)
        token_dict = user.push_notification_tokens['apns']
        for apns_token in token_dict:
            self.apns_notifications.append(Notification(payload=apns_payload, token=token_dict[apns_token]))

    def commit(self):
        """
        仅发送队列中的消息
        """
        response = apns_client.send_notification_batch(notifications=self.apns_notifications,
                                                       topic=PUSH_NOTIFICATION_CLIENT_PACKAGE_NAME_IOS)
        self.apns_notifications.clear()

        for device in response:
            if response[device] == 'BadDeviceToken':
                # TODO: remove this token
                pass

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
    return re.sub(r'[#!>_+*-]+ |[*`\[\]-]+|\d+\. |\(http.+?\)|<.+?>|\s', '', content)
