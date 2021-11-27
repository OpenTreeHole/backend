import time
from datetime import datetime, timezone, timedelta

from django.conf import settings
from django.contrib.auth import get_user_model
from django.core.cache import cache
from django.utils.dateparse import parse_datetime
from rest_framework.authtoken.models import Token
from rest_framework.test import APITestCase

from api.models import Division, Tag, Hole, Floor, Report, Message

User = get_user_model()

USERNAME = 'my email'
PASSWORD = 'my password'
EMAIL = 'test@test.com'
VERY_LONG_TIME = settings.VERY_LONG_TIME
CONTENT = 'This is a content'


def basic_setup(self):
    admin = User.objects.create_superuser('admin')
    admin.nickname = 'admin nickname'
    admin.save()
    user = User.objects.create_user(email=USERNAME, password=PASSWORD)

    self.admin = admin
    self.user = user

    self.client.credentials(HTTP_AUTHORIZATION='Token ' + Token.objects.get(user=user).key)

    division, created = Division.objects.get_or_create(name='树洞')
    for tag_name in ['tag A1', 'tag A2', 'tag B1', 'tag B2']:
        Tag.objects.create(name=tag_name, temperature=0)
    for i in range(6):
        hole = Hole.objects.create(division=division, reply=0, mapping={1: 'Jack'})
        tag_names = ['tag A1', 'tag A2'] if i % 2 == 0 else ['tag B1', 'tag B2']
        tags = Tag.objects.filter(name__in=tag_names)
        hole.tags.set(tags)
        for j in range(6):
            Floor.objects.create(
                hole=hole, anonyname='Jack', user=user,
                content='**Hole#{}; Floor No.{}**'.format(i + 1, j + 1)
            )
    return {
        'admin': admin,
        'admin_token': Token.objects.get(user=admin).key,
        'user': user,
        'user_token': Token.objects.get(user=user).key
    }


class PermissionTests(APITestCase):
    def setUp(self):
        r = basic_setup(self)
        self.user = r.get('user')
        self.user_token = r.get('user_token')
        self.admin = r.get('admin')
        self.admin_token = r.get('admin_token')

    def test_not_authenticated(self):
        self.client.credentials(HTTP_AUTHORIZATION='')
        for method in ['get', 'post', 'put', 'delete']:
            for url in ['/holes', '/floors', '/tags', '/user/favorites', '/reports']:
                loc = locals()
                exec('r = self.client.{method}("{url}")'.format(method=method, url=url), globals(), loc)
                r = loc['r']
                self.assertEqual(r.status_code, 401)

    def test_another_user(self):
        another_user = User.objects.create_user('another user')
        another_user_token = Token.objects.get(user=another_user).key
        self.client.credentials(HTTP_AUTHORIZATION='Token ' + another_user_token)

        r = self.client.put('/holes/1')
        self.assertEqual(r.status_code, 403)

        r = self.client.delete('/holes/1')
        self.assertEqual(r.status_code, 403)

        # r = self.client.put('/floors/1')
        # self.assertEqual(r.status_code, 403)

        r = self.client.delete('/floors/1')
        self.assertEqual(r.status_code, 403)

    def test_admin(self):
        self.client.credentials(HTTP_AUTHORIZATION='Token ' + self.admin_token)

        r = self.client.put('/holes/1')
        self.assertEqual(r.status_code, 200)

        r = self.client.delete('/holes/1')
        self.assertEqual(r.status_code, 405)

        r = self.client.put('/floors/1')
        self.assertEqual(r.status_code, 200)

        r = self.client.delete('/floors/1')
        self.assertEqual(r.status_code, 204)

    def test_silent(self):
        silent_user = User.objects.create_user('silented user')
        silent_user.permission['silent'][1] = VERY_LONG_TIME
        silent_user.save()
        silented_user_token = Token.objects.get(user=silent_user).key
        self.client.credentials(HTTP_AUTHORIZATION='Token ' + silented_user_token)

        data = {
            'content': CONTENT,
            'division_id': 1,
            'hole_id': 1,
            'tag_names': ['tag'],
        }
        r = self.client.post('/holes', data)
        self.assertEqual(r.status_code, 403)

        r = self.client.post('/floors', data)
        self.assertEqual(r.status_code, 403)

    def test_tags(self):
        self.client.credentials(HTTP_AUTHORIZATION='Token ' + self.user_token)
        r = self.client.get('/tags')
        self.assertEqual(r.status_code, 200)
        r = self.client.post('/tags')
        self.assertEqual(r.status_code, 403)
        r = self.client.put('/tags')
        self.assertEqual(r.status_code, 403)
        r = self.client.delete('/tags')
        self.assertEqual(r.status_code, 403)

    def test_reports(self):
        r = self.client.get('/reports')
        self.assertEqual(r.status_code, 403)
        r = self.client.put('/reports')
        self.assertEqual(r.status_code, 403)
        r = self.client.delete('/reports')
        self.assertEqual(r.status_code, 403)


class LoginLogoutTests(APITestCase):
    email = EMAIL
    password = "iasjludfnbasvdfljnhk"
    wrong_password = "saasor;lkjjhgny"

    def setUp(self):
        self.user = User.objects.create_user(email=self.email, password=self.password)

    def test_login(self):
        # 正确密码
        r = self.client.post("/login", {
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

    def test_logout(self):
        token = Token.objects.get(user=self.user).key
        self.client.credentials(HTTP_AUTHORIZATION=f'Token {token}')
        r = self.client.get('/logout')
        self.assertEqual(r.status_code, 200)
        self.assertNotEqual(Token.objects.get(user=self.user).key, token)


class RegisterTests(APITestCase):
    email = EMAIL
    another_email = 'another@test.com'
    wrong_email = "test@foo.com"
    password = "fsdvkhjng"
    simple_password = '123456'
    new_password = 'jwhkerbb4v5'
    verification = None

    def setUp(self):
        User.objects.create_user(self.another_email, password=self.password)

    def register_verify(self):
        # 正确校验
        r = self.client.get("/verify/email", {"email": self.email})
        self.assertEqual(r.status_code, 200)
        self.assertIsNotNone(r.json()['message'])
        self.assertIsNotNone(cache.get(self.email))
        self.verification = cache.get(self.email)

        # 错误域名
        r = self.client.get("/verify/email", {"email": self.wrong_email})
        self.assertEqual(r.status_code, 400)
        self.assertEqual('邮箱不在白名单内！', r.data['message'])
        self.assertIsNone(cache.get(self.wrong_email))

        # # 重复邮箱
        # r = self.client.get("/verify/email", {"email": self.another_email})
        # self.assertEqual(r.status_code, 400)
        # self.assertEqual("该用户已注册！", r.data['message'])
        # self.assertIsNone(cache.get(self.another_email))

    def register(self):
        expected_users = User.objects.count() + 1
        # 正确注册
        r = self.client.post("/register", {
            "email": self.email,
            "password": self.password,
            'verification': self.verification
        })
        self.assertEqual(r.status_code, 201)
        self.assertIsNotNone(r.data['message'])
        user = User.objects.get(email=self.email)
        Token.objects.get(user=user)

        # 重复注册
        r = self.client.post("/register", {
            "email": self.email,
            "password": self.password,
            'verification': self.verification
        })
        self.assertEqual(r.status_code, 400)
        # self.assertEqual(r.json(), {"message": "该用户已注册！"})
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
            "email": self.email,
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

    def modify_password(self):
        r = self.client.put("/register", {
            "email": self.email,
            "password": self.new_password,
            'verification': self.verification
        })
        self.assertEqual(r.status_code, 200)
        self.assertEqual(r.json(), {'message': '已重置密码'})

        r = self.client.post('/login', {
            'email': EMAIL,
            'password': self.new_password,
        })
        self.assertEqual(r.status_code, 200)

    def test(self):
        self.register_verify()
        self.register()
        self.modify_password()


class DivisionTests(APITestCase):
    def setUp(self):
        self.admin = None
        self.user = None
        basic_setup(self)

    def test_get(self):
        r = self.client.get('/divisions/1')
        self.assertEqual(r.status_code, 200)

        r = self.client.get('/divisions')
        self.assertEqual(r.status_code, 200)

    def test_put(self):
        r = self.client.put('/divisions/1')
        self.assertEqual(r.status_code, 403)

        self.client.force_authenticate(user=self.admin)
        r = self.client.put('/divisions/1', {
            'name': 'name',
            'description': 'description',
            'pinned': [1, 2]
        })
        self.assertEqual(r.status_code, 200)
        self.assertEqual(r.json()['name'], 'name')
        self.assertEqual(r.json()['description'], 'description')
        self.assertEqual(len(r.json()['pinned']), 2)

    def test_order(self):
        division = Division.objects.create(name='name', pinned=[4, 2, 5])
        r = self.client.get(f'/divisions/{division.id}')
        ids = list(map(lambda hole: hole['hole_id'], r.json()['pinned']))
        self.assertEqual(ids, [4, 2, 5])


class HoleTests(APITestCase):
    content = 'This is a content'

    def setUp(self):
        self.admin = None
        self.user = None
        basic_setup(self)

    def test_post(self):
        r = self.client.post('/holes', {
            'content': self.content,
            'division_id': 1,
            'tags': [{'name': 'tag1', 'color': 'red'}, {'name': 'tag2', 'color': 'blue'}, {'name': 'tag3', 'color': 'green'}]
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
            'length': 3,
        })
        self.assertEqual(r.status_code, 200)
        self.assertEqual(len(r.data), 3)

    def test_get_by_tag(self):
        r = self.client.get('/holes', {
            'start_time': datetime.now(timezone.utc).isoformat(),
            'length': 3,
            'tag': 'tag A1'
        })
        self.assertEqual(r.status_code, 200)
        self.assertEqual(len(r.data), 3)

    def test_get_one(self):
        r = self.client.get('/holes/1')
        self.assertEqual(r.status_code, 200)

    def test_put(self):
        self.client.force_authenticate(user=self.admin)
        r = self.client.put('/holes/1', {
            'view': 2,
            'tags': [{'name': 'tag A1', 'color': 'red'}, {'name': 'tag B1', 'color': 'red'}]
        })
        self.client.force_authenticate(user=self.user)
        self.assertEqual(r.status_code, 200)
        hole = Hole.objects.get(pk=1)
        self.assertEqual(hole.view, 2)
        tags = set(hole.tags.values_list('name', 'color'))
        self.assertEqual(tags, {('tag A1', 'red'), ('tag B1', 'red')})


class FloorTests(APITestCase):

    def setUp(self):
        self.admin = None
        self.user = None
        basic_setup(self)

    def test_post(self):
        hole = Hole.objects.get(pk=1)
        mention = list(hole.floor_set.order_by('id')[:2].values_list('id', flat=True))
        r = self.client.post('/floors', {
            'content': CONTENT,
            'hole_id': 1,
            'mention': mention,
        })
        self.assertEqual(r.status_code, 201)
        self.assertEqual(r.data['message'], '发表成功！')
        floor = Floor.objects.get(content=CONTENT)
        self.assertEqual(list(floor.mention.values_list('id', flat=True)), mention)

    def test_get(self):
        r = self.client.get('/floors', {
            'hole_id': 1,
            'start_floor': 2,
            'length': 4,
        })
        self.assertEqual(r.status_code, 200)
        self.assertEqual(len(r.json()), 4)
        self.assertEqual(r.json()[0]['hole_id'], 1)
        self.assertEqual(r.json()[0]['is_me'], True)

    def test_get_one(self):
        r = self.client.get('/floors/1')
        self.assertEqual(r.status_code, 200)
        self.assertEqual(r.json()['floor_id'], 1)

    def test_search(self):
        r = self.client.get('/floors', {
            's': 'no.2'
        })
        self.assertEqual(r.status_code, 200)
        self.assertEqual(len(r.json()), 6)
        self.assertEqual('**Hole#6; Floor No.2**', r.json()[0]['content'])

    def test_wrong_search(self):
        r = self.client.get('/floors', {
            'hole_id': 1,
            's': '*'
        })
        self.assertEqual(r.status_code, 200)
        self.assertEqual(len(r.json()), 0)

    def test_put_anonyname(self):
        self.client.force_authenticate(user=self.admin)
        r = self.client.put('/floors/1', {
            'anonyname': 'anonyname'
        })
        floor = Floor.objects.get(pk=1)
        self.assertEqual(r.status_code, 200)
        self.assertEqual(r.json()['anonyname'], 'anonyname')
        self.assertEqual(floor.anonyname, 'anonyname')

    def test_put(self):
        original_content = Floor.objects.get(pk=1).content
        r = self.client.put('/floors/1', {
            'content': 'Modified',
            'like': 'add',
            'fold': ['fold1', 'fold2'],
            'mention': [1, 2, 3]
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
        self.assertEqual(floor.fold, ['fold1', 'fold2'])
        self.assertEqual(list(floor.mention.values_list('id', flat=True)), [1, 2, 3])
        # 取消点赞
        r = self.client.put('/floors/1', {'like': 'cancel'})
        self.assertEqual(r.json()['like'], 0)
        self.assertEqual(r.json()['liked'], False)

    def test_delete(self):
        original_content = Floor.objects.get(pk=2).content
        r = self.client.delete('/floors/2')
        floor = Floor.objects.get(pk=2)
        self.assertEqual(r.status_code, 204)
        self.assertEqual(r.data['content'], '该内容已被作者删除')
        self.assertEqual(Floor.objects.get(pk=2).deleted, True)
        self.assertEqual(floor.history[0]['altered_by'], self.user.pk)
        self.assertEqual(floor.history[0]['content'], original_content)
        # 测试管理员删除
        self.client.force_authenticate(user=self.admin)
        r = self.client.delete('/floors/2')
        self.assertEqual(r.data['content'], '该内容因违反社区规范被删除')
        r = self.client.delete('/floors/2', {'delete_reason': 'reason'})
        self.assertEqual(r.data['content'], 'reason')
        self.client.force_authenticate(user=self.user)


class TagTests(APITestCase):
    def setUp(self):
        self.admin = None
        self.user = None
        basic_setup(self)
        Tag.objects.filter(name='tag B1').update(temperature=1)

    def test_get(self):
        r = self.client.get('/tags')
        self.assertEqual(r.status_code, 200)
        self.assertEqual(len(r.json()), 4)
        for tag in r.json():
            if tag['name'] == 'tag B2':
                self.assertEqual(tag['temperature'], 3)

    def test_search(self):
        r = self.client.get('/tags', {'s': 'b'})
        self.assertEqual(r.status_code, 200)
        self.assertEqual(len(r.json()), 2)
        for tag in r.json():
            self.assertIn('B', tag['name'])
        self.assertEqual(r.json()[1]['temperature'], 1)

    def test_post(self):
        self.client.force_authenticate(user=self.admin)
        # 正确提交
        r = self.client.post('/tags', {'name': 'new tag'})
        self.assertEqual(r.status_code, 201)
        self.assertEqual(r.json()['name'], 'new tag')
        self.assertEqual(r.json()['temperature'], 0)
        Tag.objects.get(name='new tag')
        # 名称过长
        r = self.client.post('/tags', {'name': ' '.join(str(i) for i in range(settings.MAX_TAG_LENGTH))})
        self.assertEqual(r.status_code, 400)

    def test_put(self):
        self.client.force_authenticate(user=self.admin)
        r = self.client.put('/tags/1', {'name': 'new name', 'temperature': 42})
        self.assertEqual(r.status_code, 200)
        self.assertEqual(r.json()['name'], 'new name')
        self.assertEqual(r.json()['temperature'], 42)

    def test_delete(self):
        self.client.force_authenticate(user=self.admin)
        pk = Tag.objects.create(name='delete').pk
        r = self.client.delete('/tags/{}'.format(pk))
        self.assertEqual(r.status_code, 204)
        self.assertEqual(r.data, None)
        self.assertFalse(Tag.objects.filter(pk=pk).exists())


class UserTests(APITestCase):
    def setUp(self):
        self.admin = None
        self.user = None
        basic_setup(self)

    def test_get(self):
        r = self.client.get('/users')
        self.assertEqual(r.status_code, 200)

    def test_put(self):
        self.client.force_authenticate(user=self.admin)
        config = {'show_folded': 'show', 'notify': ['reply', 'favorite']}
        permission = {'admin': '2000-01-01T00:00:00+00:00', 'silent': {}}
        r = self.client.put('/users/2', {
            'favorites': [1, 2],
            'config': config,
            'permission': permission
        })
        self.assertEqual(r.status_code, 200)
        self.assertEqual(r.json()['config'], config)
        self.assertEqual(r.json()['permission'], permission)

    def test_favorites(self):
        def post():
            r = self.client.post('/user/favorites', {'hole_id': 1})
            self.assertEqual(r.status_code, 201)
            self.assertEqual(r.json(), {'message': '收藏成功', 'data': [1]})
            self.assertEqual(User.objects.get(email=USERNAME).favorites.filter(pk=1).exists(), True)

        def get():
            r = self.client.get('/user/favorites')
            self.assertEqual(r.status_code, 200)
            self.assertEqual(len(r.json()), 2)

        def put():
            r = self.client.put('/user/favorites', {'hole_ids': [2, 3]})
            self.assertEqual(r.status_code, 200)
            self.assertEqual(r.json(), {'message': '修改成功', 'data': [2, 3]})
            ids = User.objects.get(email=USERNAME).favorites.values_list('id', flat=True)
            self.assertEqual([2, 3], list(ids))

        def delete():
            r = self.client.delete('/user/favorites', {'hole_id': 2})
            self.assertEqual(r.status_code, 200)
            self.assertEqual(r.json(), {'message': '删除成功', 'data': [3]})
            ids = User.objects.get(email=USERNAME).favorites.values_list('id', flat=True)
            self.assertEqual([3], list(ids))

        post()
        put()
        get()
        delete()


class ReportTests(APITestCase):
    def setUp(self):
        self.admin = None
        self.user = None
        basic_setup(self)
        Report.objects.create(hole_id=1, floor_id=1, reason='default', dealed=False)
        Report.objects.create(hole_id=1, floor_id=2, reason='default', dealed=False)
        Report.objects.create(hole_id=1, floor_id=3, reason='default', dealed=True)
        Report.objects.create(hole_id=1, floor_id=4, reason='default', dealed=True)

    def test_post(self):
        r = self.client.post('/reports', {'floor_id': 5, 'reason': 'report floor 1'})
        self.assertEqual(r.status_code, 201)
        self.assertIsNotNone(r.json()['floor'])
        self.assertEqual(r.json()['reason'], 'report floor 1')
        self.assertTrue(Report.objects.filter(reason='report floor 1').exists())

    def test_get(self):
        self.client.force_authenticate(user=self.admin)
        r = self.client.get('/reports', {'category': 'not_dealed'})
        self.assertEqual(r.status_code, 200)
        self.assertEqual(len(r.json()), 2)
        for report in r.json():
            self.assertTrue(report['report_id'] == 1 or report['report_id'] == 2)

        r = self.client.get('/reports')
        self.assertEqual(r.status_code, 200)
        self.assertEqual(len(r.json()), 2)
        for report in r.json():
            self.assertTrue(report['report_id'] == 1 or report['report_id'] == 2)

        r = self.client.get('/reports', {'category': 'dealed'})
        self.assertEqual(r.status_code, 200)
        self.assertEqual(len(r.json()), 2)
        for report in r.json():
            self.assertTrue(report['report_id'] == 3 or report['report_id'] == 4)

        r = self.client.get('/reports', {'category': 'all'})
        self.assertEqual(r.status_code, 200)
        self.assertEqual(len(r.json()), 4)

    def test_get_one(self):
        self.client.force_authenticate(user=self.admin)
        r = self.client.get('/reports/1')
        self.assertEqual(r.status_code, 200)

    def test_delete(self):
        self.client.force_authenticate(user=self.admin)
        r = self.client.delete('/reports/1', {
            'deal': {
                'fold': ['fold 1', 'fold 2'],
                'delete': 'test delete',
                'silent': 3,
            }
        })
        self.assertEqual(r.status_code, 204)
        floor = Floor.objects.get(pk=1)
        self.assertEqual(floor.fold, ['fold 1', 'fold 2'])
        self.assertEqual(floor.deleted, True)
        self.assertEqual(floor.content, 'test delete')
        user = User.objects.get(email=USERNAME)
        self.assertTrue(parse_datetime(user.permission['silent']['1']) - datetime.now(timezone.utc) < timedelta(days=3, minutes=1))
        r = self.client.get('/reports/1')
        self.assertEqual(r.json()['dealed_by'], self.admin.nickname)


class MessageTests(APITestCase):
    def setUp(self):
        self.admin = None
        self.user = None
        basic_setup(self)
        Message.objects.create(message='message', user=self.user, has_read=True)
        Message.objects.create(message='message', user=self.user, has_read=False)
        Message.objects.create(message='message', user=self.user, has_read=True)

    def test_post_share_email(self):
        r = self.client.post('/messages', {
            'to': 1,
            'share_email': True
        })
        self.assertEqual(r.status_code, 201)
        self.assertTrue(Message.objects.filter(message__contains='邮箱').exists())
        self.assertIsNotNone(r.json()['message'])

    def test_post_send_notification(self):
        self.client.force_authenticate(user=self.admin)
        r = self.client.post('/messages', {
            'to': 1,
            'message': 'hi'
        })
        self.assertEqual(r.status_code, 201)
        self.assertTrue(Message.objects.filter(message__contains='hi').exists())
        self.assertIsNotNone(r.json()['message'])

    def test_get_one(self):
        r = self.client.get('/messages/1')
        self.assertEqual(r.status_code, 200)
        self.assertEqual(r.json()['message'], 'message')

    def test_get_many(self):
        r = self.client.get('/messages', {
            'not_read': True,
            'start_time': '1970-01-01T00:00:00+00:00'
        })
        self.assertEqual(r.status_code, 200)
        self.assertEqual(len(r.json()), 1)

    def test_put(self):
        self.client.force_authenticate(user=self.admin)
        r = self.client.put('/messages/3', {
            'message': 'new',
            'has_read': True
        })
        self.assertEqual(r.status_code, 200)
        self.assertEqual(r.json()['message'], 'new')
        self.assertEqual(r.json()['has_read'], True)
