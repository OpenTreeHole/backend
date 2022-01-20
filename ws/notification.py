import uuid

from channels.db import database_sync_to_async

from api.models import Message
from api.serializers import MessageSerializer
from ws.utils import MyJsonWebsocketConsumer


class NotificationConsumer(MyJsonWebsocketConsumer):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.user = None

    async def connect(self):
        await self.accept()
        self.user = self.scope["user"]
        if self.user.is_authenticated:
            await self.channel_layer.group_add(f'user-{self.user.id}', self.channel_name)
            await self.send_json({'message': '未读消息'})
            for message in await get(self.user):
                await self.send_json(message)
        else:
            uid = str(uuid.uuid4())
            await self.send_json({'uuid': uid})
            await self.channel_layer.group_add(uid, self.channel_name)

    async def receive_json(self, content, **kwargs):
        if not self.user.is_authenticated:
            return

        action = content.get('action', '')
        id = content.get('id')
        unread = content.get('unread', True)

        if action == 'get':
            data = await get(self.user, id, unread)
            await self.send_json(data)
        elif action == 'read' and id:
            data = await read(self.user, id, has_read=True)
            await self.send_json(data)
        elif action == 'unread' and id:
            data = await read(self.user, id, has_read=False)
            await self.send_json(data)
        elif action == 'clear':
            await clear(self.user)
            await self.send_json({'message': '所有未读消息已清空'})
        else:
            await self.send_json({'message': 'action 字段不合法'})

    async def on_send(self, event):
        await self.send_json(event['content'])


@database_sync_to_async
def get(user, id=None, unread=True):
    messages = Message.objects.filter(user=user)
    if id:
        messages = messages.filter(pk=id)
        if len(messages) == 0:
            return
        serializer = MessageSerializer(messages[0])
    else:
        if unread:
            messages = messages.filter(has_read=False)
        serializer = MessageSerializer(messages, many=True)
    return serializer.data


@database_sync_to_async
def read(user, id, has_read=True):
    messages = Message.objects.filter(user=user, pk=id)
    if len(messages) == 0:
        return
    message = messages[0]
    message.has_read = has_read
    message.save()
    serializer = MessageSerializer(message)
    return serializer.data


@database_sync_to_async
def clear(user):
    Message.objects.filter(user=user).update(has_read=True)
