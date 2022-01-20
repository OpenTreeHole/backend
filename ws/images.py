import base64
import uuid
from datetime import datetime

import httpx
import magic
from django.conf import settings

from ws.utils import MyWebsocketConsumer


class ImageConsumer(MyWebsocketConsumer):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)

    async def connect(self):
        if self.scope["user"].is_authenticated:
            await self.accept()
        else:
            await self.close()

    async def receive(self, text_data=None, bytes_data=None):
        image = bytes_data
        if not image:
            await self.send_json({'message': '内容不能为空', 'status': 'error'})
        if len(image) > settings.MAX_IMAGE_SIZE * 1024 * 1024:
            await self.send_json({'message': f'图片大小不能超过 {settings.MAX_IMAGE_SIZE} MB', 'status': 'error'})
        mime = magic.from_buffer(image, mime=True)
        if mime.split('/')[0] != 'image':
            return self.send_json({'message': '请上传图片格式', 'status': 'error'})

        await self.send_json({'message': '处理中', 'status': 'info'})
        # 上传图片
        if settings.IMAGE_BACKEND == 'github':
            print(settings.HTTP_PROXY)
            date_str = datetime.now().strftime('%Y-%m-%d')
            uid = uuid.uuid4()
            filetype = mime.split('/')[1]
            async with httpx.AsyncClient(timeout=20, proxies=settings.HTTP_PROXY) as client:
                r = await client.put(
                    url=f'https://api.github.com/repos/{settings.GITHUB_OWENER}/{settings.GITHUB_REPO}/contents/{date_str}/{uid}.{filetype}',
                    headers={'Authorization': f'token {settings.GITHUB_TOKEN}'},
                    json={
                        'content': base64.b64encode(image).decode('utf-8'),
                        'message': 'upload image',
                        'branch': settings.GITHUB_BRANCH,
                    })
            result_url = f'https://cdn.jsdelivr.net/gh/{settings.GITHUB_OWENER}/{settings.GITHUB_REPO}@{settings.GITHUB_BRANCH}/{date_str}/{uid}.{filetype}'
            if r.status_code == 201:
                await self.send_json({'message': '上传成功', 'url': result_url, 'status': 'success'})
            else:
                await self.send_json({'message': '上传失败', 'data': r.json(), 'status': 'error'})
        elif settings.IMAGE_BACKEND == 'chevereto':
            async with httpx.AsyncClient(timeout=20, proxies=settings.HTTP_PROXY) as client:
                r = await client.post(
                    url=settings.CHEVERETO_URL,
                    files={'source': image},
                    data={'key': settings.CHEVERETO_TOKEN})
            if r.status_code == 200:
                r = r.json()['image']
                await self.send_json({
                    'message': '上传成功',
                    'status': 'success',
                    'url': r['url'],
                    'medium': r.get('medium', {}).get('url', r['url']),
                    'thumb': r.get('thumb', {}).get('url', r['url'])
                })
            else:
                try:
                    message = r.json()['error']['message']
                except:
                    message = '上传失败'
                await self.send_json({'message': message, 'status': 'error'})


        else:
            await self.send_json({'message': '暂不支持图片上传', 'status': 'error'})
