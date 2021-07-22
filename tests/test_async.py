from django.contrib.auth.models import User
from rest_framework.test import APITestCase


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
