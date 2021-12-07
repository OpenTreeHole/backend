import json

from channels.generic.websocket import AsyncJsonWebsocketConsumer


class MyJsonWebsocketConsumer(AsyncJsonWebsocketConsumer):
    async def send_json(self, content, close=False):
        """
        unicode 编码 json 并发给客户端
        """
        await super().send(text_data=json.dumps(content, ensure_ascii=False), close=close)

    async def on_send(self, event):
        await self.send_json(event['content'])
