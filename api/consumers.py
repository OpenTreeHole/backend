from asgiref.sync import async_to_sync
from channels.db import database_sync_to_async
from channels.generic.websocket import AsyncJsonWebsocketConsumer
from channels.layers import get_channel_layer

from api.models import Message
from api.serializers import MessageSerializer


@database_sync_to_async
def get_unread_messages(user):
    messages = Message.objects.filter(user=user, has_read=False)
    serializer = MessageSerializer(messages, many=True)
    return serializer.data


def send_message_to_user(user_id, content):
    channel_layer = get_channel_layer()
    async_to_sync(channel_layer.group_send)(
        f'user-{user_id}',  # Channels 组名称
        {
            "type": "notification",
            "content": content,
        }
    )


class NotificationConsumer(AsyncJsonWebsocketConsumer):
    async def connect(self):
        user = self.scope["user"]
        # 仅允许已登录用户
        if user.is_authenticated:
            await self.accept()
            await self.channel_layer.group_add(f'user-{user.id}', self.channel_name)
            await self.send_json({
                'message': '未读消息',
                'data': await get_unread_messages(user)
            })
        else:
            await self.close()

    async def disconnect(self, close_code):
        pass

    async def receive_json(self, content, **kwargs):
        message = content.get('message', 'hi')

        await self.send_json({
            'message': message,
        })

    async def notification(self, event):
        await self.send_json(event['content'])
