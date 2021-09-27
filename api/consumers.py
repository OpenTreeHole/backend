from channels.generic.websocket import AsyncJsonWebsocketConsumer


class NotificationConsumer(AsyncJsonWebsocketConsumer):
    async def connect(self):
        self.user = self.scope["user"]
        print(self.user)
        await self.accept()

    async def disconnect(self, close_code):
        pass

    async def receive_json(self, content, **kwargs):
        message = content.get('message')

        await self.send_json({
            'message': message,
            'hi': 'hi',
        })
