from django.contrib.auth import get_user_model
from rest_framework.test import APITestCase

from tests.test_apis import basic_setup
from utils.auth import many_hashes

User = get_user_model()


class IndexTests(APITestCase):
    """hi"""

    def test_get(self):
        r = self.client.get("/")
        self.assertEqual(r.status_code, 200)
        self.assertEqual(r.data, {"message": "Hello world!"})


class MessageTest(APITestCase):
    def setUp(self):
        basic_setup(self)
        self.another_user = User.objects.create_user('another user')

    def test_post(self):
        r = self.client.post('/messages', {
            'from': User.objects.get(identifier=many_hashes('another user')).pk,
            'to': 1,
            'share_email': True,
        })
        self.assertEqual(r.status_code, 201)
        self.assertIn('message', r.json())


class EmailTest(APITestCase):
    def setUp(self):
        basic_setup(self)

    def test_wrong_email(self):
        r = self.client.post('/email/password', {
            'email': 'test',
            'password': 123456,
        })
        self.assertEqual(r.status_code, 400)
        self.assertIn('email', r.json())

    def test_wrong_path(self):
        r = self.client.post('/email/path', {
            'email': 'test@test.com',
            'password': 123456,
        })
        self.assertEqual(r.status_code, 404)
        self.assertIn('message', r.json())

    def test_password_email(self):
        r = self.client.post('/email/password', {
            'email': 'test@test.com',
            'password': 123456,
        })
        self.assertEqual(r.status_code, 202)
        self.assertIn('message', r.json())

    def test_no_password(self):
        r = self.client.post('/email/password', {
            'email': 'test@test.com',
        })
        self.assertEqual(r.status_code, 400)
        self.assertIn('message', r.json())
