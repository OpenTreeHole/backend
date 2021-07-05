from django.core.cache import cache
from rest_framework.test import APITestCase

from api.models import *


class IndexTests(APITestCase):
    """hi"""

    def test_get(self):
        r = self.client.get("/")
        self.assertEqual(r.data, {"message": "Hello world!"})


class RegisterTests(APITestCase):
    email = "test@test.com"
    another_email = 'another@test.com'
    wrong_email = "test@foo.com"
    password = "fsdvkhjng"
    simple_password = '123456'
    verification = None

    def setUp(self):
        User.objects.create_user(self.another_email, password=self.password)

    def register_verify(self):
        # 正确校验
        r = self.client.get("/verify/email", {"email": self.email})
        self.assertEqual(r.status_code, 200)
        self.assertIsNotNone(cache.get(self.email))
        self.verification = cache.get(self.email)

        # 错误域名
        r = self.client.get("/verify/email", {"email": self.wrong_email})
        self.assertEqual(r.status_code, 400)
        self.assertEqual('邮箱不在白名单内！', r.data['message'])
        self.assertIsNone(cache.get(self.wrong_email))

        # 重复邮箱
        r = self.client.get("/verify/email", {"email": self.another_email})
        self.assertEqual(r.status_code, 400)
        self.assertEqual("该用户已注册！", r.data['message'])
        self.assertIsNone(cache.get(self.another_email))

    def register(self):
        # 正确注册
        r = self.client.post("/register", {
            "email": self.email,
            "password": self.password,
            'verification': self.verification
        })
        self.assertEqual(r.status_code, 200)
        self.assertEqual(r.data['message'], '注册成功！')
        user = User.objects.get(username=self.email)
        Token.objects.get(user=user)
        Profile.objects.get(user=user)

        # 简单密码
        r = self.client.post("/register", {
            "email": self.email,
            "password": self.simple_password,
            'verification': self.verification,
        })
        self.assertEqual(r.status_code, 400)
        self.assertIn('密码', r.data['message'])
        self.assertEqual(User.objects.count(), 2)

        # 错误邮箱
        r = self.client.post("/register", {
            "email": self.another_email,
            "password": self.password,
            'verification': self.verification,
        })
        self.assertEqual(r.status_code, 400)
        self.assertEqual('注册校验未通过！', r.data['message'])
        self.assertEqual(User.objects.count(), 2)

    def test(self):
        self.register_verify()
        self.register()


class LoginTests(APITestCase):
    email = "test@test.com"
    password = "iasjludfnbasvdfljnhk"
    wrong_password = "saasor;lkjjhgny"

    def setUp(self):
        User.objects.create_user(username=self.email, password=self.password)

    def test_login(self):
        r = self.client.post(
            "/login",
            data={
                "email": self.email,
                "password": self.password,
            },
        )
        self.assertEqual(r.status_code, 200)
        self.assertIn("token", r.data)

        r = self.client.post(
            "/login",
            data={
                "email": self.email,
                "password": self.wrong_password,
            },
        )
        self.assertEqual(r.status_code, 401)
