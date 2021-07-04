from django.contrib.auth.models import User
from api.models import *

from rest_framework.test import APITestCase
from django.contrib.auth.hashers import make_password, check_password


class IndexTests(APITestCase):
    """hi"""

    def test_get(self):
        print("1")
        r = self.client.get("/")
        self.assertEqual(r.data, {"message": "Hello world!"})


class AuthenticationTests(APITestCase):
    def setUp(self):
        User.objects.create(username="username", password=make_password("password"))

    def test_login(self):
        print("2")
        r = self.client.post(
            "/login",
            data={
                "username": "username",
                "password": "password",
            },
        )
        self.assertEqual(r.status_code, 200)
        self.assertIn("token", r.data)

    def test_wrong_login(self):
        print("3")
        r = self.client.post(
            "/login",
            data={
                "username": "username",
                "password": "password1",
            },
        )
        self.assertEqual(r.status_code, 401)
