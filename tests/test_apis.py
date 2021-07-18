from datetime import datetime, timezone
import time
from django.core.cache import cache
from rest_framework.test import APITestCase
from api.models import *

USERNAME = 'my username'
PASSWORD = 'my password'
EMAIL = 'test@test.com'


def basic_setup(self):
    User.objects.create_user('admin')
    user = User.objects.create_user(username=USERNAME, password=PASSWORD)
    token = Token.objects.get(user=user)
    self.client.credentials(HTTP_AUTHORIZATION='Token ' + token.key)

    division, created = Division.objects.get_or_create(name='树洞')
    for tag_name in ['tag A1', 'tag A2', 'tag B1', 'tag B2']:
        Tag.objects.create(name=tag_name, temperature=5)
    for i in range(10):
        hole = Hole.objects.create(division=division, reply=0, mapping={1: 'Jack'})
        tag_names = ['tag A1', 'tag A2'] if i % 2 == 0 else ['tag B1', 'tag B2']
        tags = Tag.objects.filter(name__in=tag_names)
        hole.tags.set(tags)
        for j in range(10):
            Floor.objects.create(
                hole=hole, anonyname='Jack', user=user,
                content='**Hole#{}; Floor No.{}**'.format(i + 1, j + 1),
                shadow_text='Hole#{}; Floor No.{}'.format(i + 1, j + 1),
            )


class IndexTests(APITestCase):
    """hi"""

    def test_get(self):
        r = self.client.get("/")
        self.assertEqual(r.status_code, 200)
        self.assertEqual(r.data, {"message": "Hello world!"})


class LoginTests(APITestCase):
    email = EMAIL
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
    email = EMAIL
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
        expected_users = User.objects.count() + 1
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

        # 重复注册
        r = self.client.post("/register", {
            "email": self.email,
            "password": self.password,
            'verification': self.verification
        })
        self.assertEqual(r.status_code, 400)
        self.assertEqual(r.json(), {"message": "该用户已注册！"})
        self.assertEqual(User.objects.count(), expected_users)

        # 未提供验证码
        r = self.client.post("/register", {
            "email": self.email,
            "password": self.password,
            # 'verification': self.verification
        })
        self.assertEqual(r.status_code, 400)
        self.assertEqual(r.data['message'], '验证码不能为空！')
        self.assertEqual(User.objects.count(), expected_users)

        # 简单密码
        r = self.client.post("/register", {
            "email": self.wrong_email,
            "password": self.simple_password,
            'verification': self.verification,
        })
        self.assertEqual(r.status_code, 400)
        self.assertIn('密码', r.data['message'])
        self.assertEqual(User.objects.count(), expected_users)

        # 错误邮箱
        r = self.client.post("/register", {
            "email": self.wrong_email,
            "password": self.password,
            'verification': self.verification,
        })
        self.assertEqual(r.status_code, 400)
        self.assertEqual('注册校验未通过！', r.data['message'])
        self.assertEqual(User.objects.count(), expected_users)

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
            'tag_names': ['tag1', 'tag2', 'tag3']
        })
        self.assertEqual(r.status_code, 201)
        self.assertEqual(r.data['message'], '发表成功！')
        floor = Floor.objects.get(content=self.content)
        hole = floor.hole
        self.assertEqual(hole.tags.count(), 3)
        for tag in hole.tags.all():
            self.assertEqual(tag.temperature, 1)

    def test_get_by_time(self):
        time.sleep(1)
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

    def test_put(self):
        r = self.client.put('/holes/1', {
            'view': 2,
            'tag_names': ['tag A1', 'tag B1']
        })
        self.assertEqual(r.status_code, 200)
        hole = Hole.objects.get(pk=1)
        self.assertEqual(hole.view, 2)
        tag_names = set()
        for tag in hole.tags.all():
            tag_names.add(tag.name)
        self.assertEqual(tag_names, {'tag A1', 'tag B1'})


class FloorTests(APITestCase):
    content = 'This is a content'

    def setUp(self):
        basic_setup(self)
        self.user = User.objects.get(username=USERNAME)

    def test_post(self):
        hole = Hole.objects.get(pk=1)
        first_floor = hole.floor_set.order_by('id')[0]
        r = self.client.post('/floors', {
            'content': self.content,
            'hole_id': 1,
            'reply_to': first_floor.pk,
        })
        self.assertEqual(r.status_code, 201)
        self.assertEqual(r.data['message'], '发表成功！')
        floor = Floor.objects.get(content=self.content)
        self.assertEqual(floor.reply_to, first_floor.pk)

    def test_get(self):
        r = self.client.get('/floors', {
            'hole_id': 1,
            'start_floor': 3,
            'length': 5,
        })
        self.assertEqual(r.status_code, 200)
        self.assertEqual(len(r.json()), 5)
        self.assertEqual(r.json()[0]['hole_id'], 1)
        self.assertEqual(r.json()[0]['is_me'], True)

    def test_search(self):
        r = self.client.get('/floors', {
            'hole_id': 1,
            's': 'no.2'
        })
        self.assertEqual(r.status_code, 200)
        self.assertEqual(len(r.json()), 1)
        self.assertEqual('**Hole#1; Floor No.2**', r.json()[0]['content'])

    def test_wrong_search(self):
        r = self.client.get('/floors', {
            'hole_id': 1,
            's': '*'
        })
        self.assertEqual(r.status_code, 200)
        self.assertEqual(len(r.json()), 0)

    def test_put(self):
        original_content = Floor.objects.get(pk=1).content
        r = self.client.put('/floors', {
            'floor_id': 1,
            'content': 'Modified',
            'like': True,
            'folded': ['folded1', 'folded2']
        })
        self.assertEqual(r.status_code, 200)
        self.assertEqual(r.json()['content'], 'Modified')
        self.assertEqual(r.json()['like'], 1)
        self.assertEqual(r.json()['liked'], True)
        floor = Floor.objects.get(pk=1)
        self.assertEqual(floor.content, 'Modified')
        self.assertEqual(floor.like, 1)
        self.assertIn(self.user.pk, floor.like_data)
        self.assertEqual(floor.history[0]['altered_by'], self.user.pk)
        self.assertEqual(floor.history[0]['content'], original_content)
        self.assertEqual(floor.folded, ['folded1', 'folded2'])

    def test_delete(self):
        original_content = Floor.objects.get(pk=2).content
        r = self.client.delete('/floors', {'floor_id': 2})
        floor = Floor.objects.get(pk=2)
        self.assertEqual(r.status_code, 204)
        # self.assertEqual(r.json()['content'], '该内容已被作者删除')
        self.assertEqual(Floor.objects.get(pk=2).deleted, True)
        self.assertEqual(floor.history[0]['altered_by'], self.user.pk)
        self.assertEqual(floor.history[0]['content'], original_content)


class PermissionTests(APITestCase):
    def setUp(self):
        admin = User.objects.create_user('admin')
        admin.profile.permission['admin'] = '9999-01-01T00:00:00+00:00'
        admin.profile.save()
        self.admin = admin
        self.admin_token = Token.objects.get(user=admin).key

        user = User.objects.create_user('user')
        self.user = user
        self.user_token = Token.objects.get(user=user).key
        # self.client.credentials(HTTP_AUTHORIZATION='Token ' + token.key)

    def test_authentication(self):
        self.client.credentials(HTTP_AUTHORIZATION='')
        for method in ['get', 'post', 'put', 'delete']:
            for url in ['/holes', '/floors']:
                loc = locals()
                exec('r = self.client.{method}("{url}")'.format(method=method, url=url), globals(), loc)
                r = loc['r']
                self.assertEqual(r.status_code, 401)

    # def test_permission(self):
