import collections
from os import environ

from apns2.client import APNsClient
from apns2.payload import Payload as APNsPayload
from asgiref.sync import async_to_sync
from channels.layers import get_channel_layer
from django.conf import settings

from api.models import Message
from api.serializers import MessageSerializer

# APNS global definition
Notification = collections.namedtuple('Notification', ['token', 'payload'])
if settings.APNS_KEY_PATH:
    apns_client = APNsClient(settings.APNS_KEY_PATH, use_sandbox=(environ.get("HOLE_ENV") != "production"),
                             use_alternative_port=settings.APNS_USE_ALTERNATIVE_PORT)
    print('APNS Client Initialized')
else:
    apns_client = None


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
    apns_user_token_record = {}

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
        if not user.is_authenticated:
            return

        if data is None:
            data = {}

        instance = Message.objects.create(user=user, message=message, data=data, code=code)
        payload = MessageSerializer(instance).data

        # WS
        self.__send_websocket_message_to_user(user, payload)

        # APNS
        if apns_client:
            apns_payload = APNsPayload(alert=message, sound="default", badge=1, thread_id=code, custom=payload)
            token_dict = user.push_notification_tokens['apns']
            for apns_device in token_dict:
                self.apns_notifications.append(Notification(payload=apns_payload, token=token_dict[apns_device]))
                self.apns_user_token_record.update({token_dict[apns_device]: user})

    def commit(self):
        """
        发送队列中的消息
        并清除过期token
        """

        def _commit():
            if apns_client:
                response = apns_client.send_notification_batch(notifications=self.apns_notifications,
                                                               topic=settings.PUSH_NOTIFICATION_CLIENT_PACKAGE_NAME_IOS)
                self.apns_notifications.clear()

                # 清除过期token
                for token in response:
                    if response[token] == 'BadDeviceToken':
                        user = self.apns_user_token_record[token]
                        for device in user.push_notification_tokens['apns']:
                            if user.push_notification_tokens['apns'][device] == token:
                                del user.push_notification_tokens['apns'][device]
                                break
                        user.save(update_fields=['push_notification_tokens'])
                self.apns_user_token_record.clear()

        try:
            _commit()
        except:
            pass
