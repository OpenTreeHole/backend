import collections
import json
import urllib.parse

import requests
from apns2.client import APNsClient
from apns2.payload import Payload as APNsPayload
from apns2.payload import PayloadAlert
from celery import shared_task
from django.conf import settings

from OpenTreeHole.config import MIPUSH_APP_SECRET
from api.models import Message, PushToken
from api.serializers import MessageSerializer
# APNS global definition
from ws.utils import send_websocket_message_to_group

Notification = collections.namedtuple('Notification', ['token', 'payload'])
APNS = None
if settings.APNS_KEY_PATH:
    try:
        APNS = APNsClient(
            settings.APNS_KEY_PATH,
            use_sandbox=(settings.HOLE_ENV != "production"),
            use_alternative_port=settings.APNS_USE_ALTERNATIVE_PORT
        )
    except Exception as e:
        print("[E] An error occurred in APNS subroutine initialization", e)


@shared_task
def send_notifications(user_id: int, message: str, data=None, code=''):
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
            elif code == 'penalty':
                # Data is Penalty
                return f"被处罚分区ID：{data['division_id']}，处罚等级：{data['level']}，截止日期：{data['date']}"
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
    send_websocket_message_to_group(f'user-{user_id}', payload)
    # 发送 APNS 通知
    if APNS:
        push_tokens = PushToken.objects.filter(user_id=user_id, service='apns').values_list('token', flat=True)
        if push_tokens:
            # 准备数据
            apns_notifications = []
            unread_count = Message.objects.filter(user_id=user_id, has_read=False).count()
            apns_payload = APNsPayload(
                alert=PayloadAlert(title=instance.message, body=_generate_subtitle(data, code)),
                sound="default",
                badge=unread_count,
                thread_id=instance.code,
                custom=payload
            )
            for push_token in push_tokens:
                apns_notifications.append(Notification(payload=apns_payload, token=push_token))
            # 发送数据
            response = APNS.send_notification_batch(
                notifications=apns_notifications,
                topic=settings.PUSH_NOTIFICATION_CLIENT_PACKAGE_NAME_IOS
            )
            print('APNS Response', response)
            # 清除过期token
            for token in response:
                if response[token] == 'BadDeviceToken':
                    PushToken.objects.filter(user_id=user_id, token=token).delete()

    if MIPUSH_APP_SECRET:
        push_tokens = PushToken.objects.filter(user_id=user_id, service='mipush').values_list('token', flat=True)
        # Only send request if token is not empty
        if push_tokens:
            try:
                response_json = requests.post(
                    "https://api.xmpush.xiaomi.com/v2/message/regid",
                    headers={"Authorization": f"key={MIPUSH_APP_SECRET}"},
                    data={
                        "registration_id": ','.join(push_tokens),
                        "restricted_package_name": settings.PUSH_NOTIFICATION_CLIENT_PACKAGE_NAME_ANDROID,
                        "title": instance.message,
                        "description": _generate_subtitle(data, code),
                        "payload": urllib.parse.urlencode({"data": json.dumps(data, ensure_ascii=False), "code": code}),
                        "extra.notify_effect": '1'
                    }).json()
                print('MiPush Response', response_json)
                # 清除过期token
                try:
                    bad_ids = response_json['data']['bad_regids']
                    if bad_ids:
                        for bad_id in bad_ids.split(','):
                            PushToken.objects.filter(user_id=user_id, token=bad_id).delete()
                except KeyError:
                    pass
            except Exception as e:
                print("[E] An error occurred in MiPush subroutine:", e)
