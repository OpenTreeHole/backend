import collections

from apns2.client import APNsClient
from apns2.payload import Payload as APNsPayload
from apns2.payload import PayloadAlert
from asgiref.sync import async_to_sync
from celery import shared_task
from channels.layers import get_channel_layer
from django.conf import settings

from api.models import Message
from api.serializers import MessageSerializer

# APNS global definition
Notification = collections.namedtuple('Notification', ['token', 'payload'])
APNS = None
if settings.APNS_KEY_PATH:
    try:
        APNS = APNsClient(
            settings.APNS_KEY_PATH,
            use_sandbox=(settings.HOLE_ENV != "production"),
            use_alternative_port=settings.APNS_USE_ALTERNATIVE_PORT
        )
        print('APNS Client Initialized')
    except Exception as e:
        print(e)


@shared_task
def send_notifications(user_id: int, message: str, data=None, code=''):
    def _send_websocket_message_to_user(user_id: int, content: dict):
        """
        向用户发送 WebSocket 消息
        Args:
            user_id: 用户 id
            content: 消息内容

        Returns: None
        """
        channel_layer = get_channel_layer()
        async_to_sync(channel_layer.group_send)(
            f'user-{user_id}',  # Channels 组名称
            {
                "type": "on_send",
                "content": content,
            }
        )

    def _generate_subtitle(data, code: str):
        """
        生成消息的副标题
        Args:
            data: 消息数据
            code: 消息类型

        Returns: 副标题
        """
        try:
            if code == 'mention' or code == 'favorite' or code == 'modify':
                # Data is Floor
                return f"{data['anonyname']}：{data['content']}"
            elif code == 'report':
                # Data is Report
                return f"内容：{data['floor']['content']}，理由：{data['reason']}"
        except Exception:
            return None

    if not user_id:
        return
    # 创建对象
    if data is None:
        data = {}
    instance = Message.objects.create(user_id=user_id, message=message, data=data, code=code)
    payload = MessageSerializer(instance).data
    # 发送 websocket 通知
    _send_websocket_message_to_user(instance.user_id, payload)
    # 发送 APNS 通知
    if APNS:
        # 准备数据
        apns_notifications = []
        apns_user_token_record = {}
        user = instance.user
        apns_payload = APNsPayload(
            alert=PayloadAlert(title=instance.message, body=_generate_subtitle(data, code)),
            sound="default",
            # badge=1,
            thread_id=instance.code,
            custom=payload
        )
        token_dict = user.push_notification_tokens['apns']
        for apns_device in token_dict:
            apns_notifications.append(Notification(payload=apns_payload, token=token_dict[apns_device]))
            apns_user_token_record.update({token_dict[apns_device]: user})
        # 发送数据
        response = APNS.send_notification_batch(
            notifications=apns_notifications,
            topic=settings.PUSH_NOTIFICATION_CLIENT_PACKAGE_NAME_IOS
        )
        # 清除过期token
        for token in response:
            if response[token] == 'BadDeviceToken':
                user = apns_user_token_record[token]
                for device in user.push_notification_tokens['apns']:
                    if user.push_notification_tokens['apns'][device] == token:
                        del user.push_notification_tokens['apns'][device]
                        break
                user.save(update_fields=['push_notification_tokens'])
