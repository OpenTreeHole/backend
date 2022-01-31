import json

from asgiref.sync import async_to_sync
from channels.generic.websocket import AsyncJsonWebsocketConsumer, AsyncWebsocketConsumer
from channels.layers import get_channel_layer


class MyJsonWebsocketConsumer(AsyncJsonWebsocketConsumer):
    async def send_json(self, content, close=False):
        """
        unicode 编码 json 并发给客户端
        """
        await super().send(text_data=json.dumps(content, ensure_ascii=False), close=close)


class MyWebsocketConsumer(AsyncWebsocketConsumer):
    async def send_json(self, content, close=False):
        await super().send(text_data=json.dumps(content, ensure_ascii=False), close=close)


def send_websocket_message_to_group(group: str, content: dict):
    """
    向 uuid 发送 WebSocket 消息
    Args:
        group: Channels 组名称
        content: 消息内容

    Returns: None
    """
    channel_layer = get_channel_layer()
    async_to_sync(channel_layer.group_send)(
        group,
        {
            "type": "on_send",
            "content": content,
        }
    )


async def async_send_websocket_message_to_group(group: str, content: dict):
    channel_layer = get_channel_layer()
    await channel_layer.group_send(
        group,
        {
            "type": "on_send",
            "content": content,
        }
    )
