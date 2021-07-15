from datetime import datetime, timezone
from time import sleep

from django.core.cache import cache
from rest_framework.test import APITestCase

from api.models import *


def basic_setup(self):
    user = User.objects.create_user('a user')
    token = Token.objects.get(user=user)
    self.client.credentials(HTTP_AUTHORIZATION='Token ' + token.key)

    division = Division.objects.create(name='树洞')
    for tag_name in ['tag A1', 'tag A2', 'tag B1', 'tag B2']:
        Tag.objects.create(name=tag_name, temperature=5)
    for i in range(10):
        hole = Hole.objects.create(division=division, reply=0, mapping={1: 'Jack'})
        tag_names = ['tag A1', 'tag A2'] if i % 2 == 0 else ['tag B1', 'tag B2']
        tags = Tag.objects.filter(name__in=tag_names)
        hole.tags.set(tags)
        for j in range(2):
            Floor.objects.create(
                hole=hole, anonyname='Jack', user=user,
                content='Hole#{}; Floor No.{}'.format(i, j)
            )


class IndexTests(APITestCase):
    """hi"""

    def test_get(self):
        r = self.client.get("/")
        self.assertEqual(r.data, {"message": "Hello world!"})


class LoginTests(APITestCase):
    email = "test@test.com"
    password = "iasjludfnbasvdfljnhk"
    wrong_password = "saasor;lkjjhgny"

    def setUp(self):
        User.objects.create_user(username=self.email, password=self.password)

    def test_login(self):
        # 正确密码
        r = self.client.post(
            "/login",
            data={
                "email": self.email,
                "password": self.password,
            },
        )
        self.assertEqual(r.status_code, 200)
        self.assertIn("token", r.data)

        # 错误密码
        r = self.client.post(
            "/login",
            data={
                "email": self.email,
                "password": self.wrong_password,
            },
        )
        self.assertEqual(r.status_code, 401)


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

        # 未提供验证码
        r = self.client.post("/register", {
            "email": self.email,
            "password": self.password,
            # 'verification': self.verification
        })
        self.assertEqual(r.status_code, 400)
        self.assertEqual(r.data['message'], '验证码不能为空！')
        self.assertEqual(User.objects.count(), 2)

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


class HoleTests(APITestCase):
    content = 'This is a content'

    def setUp(self):
        basic_setup(self)

    def test_post(self):
        r = self.client.post('/holes', {
            'content': self.content,
            'division_id': 1,
            'tags': ['tag1', 'tag2', 'tag3']
        })
        self.assertEqual(r.status_code, 200)
        self.assertEqual(r.data['message'], '发表成功！')
        floor = Floor.objects.get(content=self.content)
        hole = floor.hole
        self.assertEqual(hole.tags.count(), 3)

    def test_get_by_time(self):
        sleep(1)
        r = self.client.get('/holes', {
            'start_time': datetime.now(timezone.utc).isoformat(),
            'length': 5,
        })
        self.assertEqual(r.status_code, 200)
        self.assertEqual(len(r.data), 5)

    def test_get_by_tag(self):
        r = self.client.get('/holes', {
            'start_time': datetime.now(timezone.utc).isoformat(),
            'length': 5,
            'tag': 'tag A1'
        })
        self.assertEqual(r.status_code, 200)
        self.assertEqual(len(r.data), 5)

    def test_get_one(self):
        r = self.client.get('/holes/1')
        self.assertEqual(r.status_code, 200)


class FloorTests(APITestCase):
    content = 'This is a content'

    def setUp(self):
        basic_setup(self)

    def test_post(self):
        hole = Hole.objects.get(pk=1)
        first_floor = hole.floor_set.order_by('id')[0]
        r = self.client.post('/floors', {
            'content': self.content,
            'hole_id': 1,
            'reply_to': first_floor.pk,
        })
        self.assertEqual(r.status_code, 200)
        self.assertEqual(r.data['message'], '发表成功！')
        floor = Floor.objects.get(content=self.content)
        self.assertEqual(floor.reply_to, first_floor.pk)
