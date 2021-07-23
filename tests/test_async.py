from django.contrib.auth.models import User
from rest_framework.test import APITestCase

from tests.test_apis import basic_setup


class IndexTests(APITestCase):
    """hi"""

    def test_get(self):
        r = self.client.get("/")
        self.assertEqual(r.status_code, 200)
        self.assertEqual(r.data, {"message": "Hello world!"})


class ImageTests(APITestCase):
    def setUp(self):
        self.user = User.objects.create_user(username='username')

    def test_post(self):
        self.client.force_authenticate(user=self.user)
        with open('tests/image.jpg', 'rb') as image:
            r = self.client.post('/images', {'image': image}, format='multipart')
        self.assertEqual(r.status_code, 202)


class MessageTest(APITestCase):
    def setUp(self):
        basic_setup(self)
        self.another_user = User.objects.create_user('another user')

    def test_post(self):
        r = self.client.post('/messages', {
            'from': User.objects.get(username='another user').pk,
            'to': 1,
            'share_email': True,
        })
        self.assertEqual(r.status_code, 201)
        self.assertEqual(r.json(), {'message': '已发送通知'})
