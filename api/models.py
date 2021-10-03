from datetime import datetime, timezone

from django.conf import settings
from django.contrib.auth.base_user import AbstractBaseUser, BaseUserManager
from django.db import models
from django.db.models import F
from django.db.models.signals import post_save, m2m_changed
from django.dispatch import receiver
from django.utils.dateparse import parse_datetime
from rest_framework.authtoken.models import Token

from api.utils import send_message_to_user


class Division(models.Model):
    name = models.CharField(max_length=32, unique=True)
    description = models.TextField(null=True)

    def __str__(self):
        return self.name


class Tag(models.Model):
    name = models.CharField(max_length=8, unique=True)
    temperature = models.IntegerField(db_index=True, default=0, help_text="该标签下的主题帖数")

    def __str__(self):
        return self.name


# 主题帖
class Hole(models.Model):
    time_created = models.DateTimeField(auto_now_add=True)
    time_updated = models.DateTimeField(auto_now=True, db_index=True)
    tags = models.ManyToManyField(Tag, blank=True)
    division = models.ForeignKey(Division, on_delete=models.CASCADE, help_text="分区")
    view = models.IntegerField(db_index=True, default=0, help_text="浏览量")
    reply = models.IntegerField(db_index=True, default=-1, help_text="楼层数")  # 如果只有首条帖子的话认为回复数为零
    mapping = models.JSONField(default=dict, help_text='匿名到真实用户的对应')  # {user.id: anonymous_name}

    # key_floors 首条和末条回帖，动态生成

    def __str__(self):
        return f'树洞#{self.pk}'


class Floor(models.Model):
    """
    history:
        'content': floor.content,                               # 原内容
        'altered_by': request.user.pk,                          # 修改者 id
        'altered_time': datetime.now(timezone.utc).isoformat()  # 修改时间
    """
    hole = models.ForeignKey(Hole, on_delete=models.CASCADE)
    content = models.TextField()
    shadow_text = models.TextField()  # 去除markdown关键字的文本，方便搜索
    anonyname = models.CharField(max_length=16)
    user = models.ForeignKey(settings.AUTH_USER_MODEL, models.CASCADE)
    reply_to = models.ForeignKey('self', on_delete=models.SET_NULL, null=True)
    time_created = models.DateTimeField(auto_now_add=True)
    time_updated = models.DateTimeField(auto_now=True)
    like = models.IntegerField(default=0, db_index=True)  # 赞同数
    like_data = models.JSONField(default=list)  # 点赞记录，主键列表
    deleted = models.BooleanField(default=False)  # 仅作为前端是否显示删除按钮的依据
    history = models.JSONField(default=list)  # 修改记录，字典列表
    fold = models.JSONField(default=list)  # 折叠原因，字符串列表（原因由前端提供）

    def __str__(self):
        return f"{self.content[:50]}"


class Report(models.Model):
    hole = models.ForeignKey(Hole, on_delete=models.CASCADE)
    floor = models.ForeignKey(Floor, on_delete=models.CASCADE)
    reason = models.TextField()
    time_created = models.DateTimeField(auto_now_add=True)
    time_updated = models.DateTimeField(auto_now=True)
    dealed = models.BooleanField(default=False, db_index=True)
    dealed_by = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.CASCADE, null=True)

    def __str__(self):
        return f"帖子#{self.hole.pk}，{self.reason}"


def default_permission():
    """
    silent 字典
        index：分区id （string） django的JSONField会将字典的int索引转换成str
        value：禁言解除时间
    """
    return {
        'admin': '1970-01-01T00:00:00+00:00',  # 管理员权限：到期时间
        'silent': {}  # 禁言
    }


def default_config():
    """
    show_folded: 对折叠内容的处理
        fold: 折叠
        hide: 隐藏
        show: 展示

    notify: 在以下场景时通知
        reply: 被回复时
        favorite: 收藏的主题帖有新帖时
        report: 被举报时
    另外，当用户权限发生变化或所发帖被修改时也会收到通知
    """
    return {
        'show_folded': 'fold',
        'notify': ['reply', 'favorite', 'report']
    }


class UserManager(BaseUserManager):
    def create_user(self, email, password=None, **extra_fields):
        if not email:
            raise ValueError('邮箱必须提供')
        email = self.normalize_email(email)
        user = self.model(email=email, **extra_fields)
        user.set_password(password)
        user.save()
        return user

    def create_superuser(self, email, password=None, **extra_fields):
        user = self.create_user(email, password, **extra_fields)
        user.permission['admin'] = settings.VERY_LONG_TIME
        user.save()
        return user


class User(AbstractBaseUser):
    email = models.CharField(max_length=150, unique=True)
    joined_time = models.DateTimeField(auto_now_add=True)
    nickname = models.CharField(max_length=32, blank=True)
    favorites = models.ManyToManyField(Hole, blank=True)
    permission = models.JSONField(default=default_permission)
    config = models.JSONField(default=default_config)

    objects = UserManager()
    USERNAME_FIELD = 'email'

    @property
    def is_admin(self):
        now = datetime.now(timezone.utc)
        expire_time = parse_datetime(self.permission['admin'])
        return expire_time > now

    def is_silenced(self, division_id):
        now = datetime.now(timezone.utc)
        silent = self.permission['silent']
        division = str(division_id)  # JSON 序列化会将字典的 int 索引转换成 str
        if not silent.get(division):  # 未设置禁言，返回 False
            return False
        else:
            expire_time = parse_datetime(silent.get(division))
            return expire_time > now

    def __str__(self):
        return f"用户#{self.pk}"


class Message(models.Model):
    user = models.ForeignKey(
        User, related_name="message_to", on_delete=models.CASCADE, db_index=True
    )
    content = models.TextField()
    has_read = models.BooleanField(default=False)
    time_created = models.DateTimeField(auto_now_add=True)

    def __str__(self):
        return f"-> {self.user.pk}: {self.content}"


# 自动在创建用户后创建其 Token
@receiver(post_save, sender=settings.AUTH_USER_MODEL)
def create_token(sender, instance=None, created=False, **kwargs):
    if created:
        Token.objects.create(user=instance)


# 自动修改 tag 的热度
@receiver(m2m_changed, sender=Hole.tags.through)
def modify_tag_temperature(sender, reverse, action, pk_set, **kwargs):
    if reverse is False and action == 'post_add':
        Tag.objects.filter(pk__in=pk_set).update(temperature=F('temperature') + 1)
    elif reverse is False and action == 'post_remove':
        Tag.objects.filter(pk__in=pk_set).update(temperature=F('temperature') - 1)


# 在数据库中创建一条消息并通过 websocket 发送给用户
def create_and_send_message(user, message):
    Message.objects.create(user=user, content=message)
    send_message_to_user(user, {'message': message})


# 收到回复后通知用户
@receiver(post_save, sender=Floor)
def notify_when_replied(sender, instance, created, **kwargs):
    if created and instance.reply_to:
        if 'reply' in instance.reply_to.user.config['notify']:
            message = f'你在 {instance.reply_to.hole} 的帖子 {instance.reply_to} 被回复了'
            create_and_send_message(instance.reply_to.user, message)

# 收藏的主题帖有新帖时通知用户
# @receiver(post_save, sender=Floor)
# def notify_when_favorites_updated(sender, instance, created, **kwargs):
#     if created and instance.hole:
#         if 'reply' in instance.reply_to.user.config['notify']:
#             message = f'你在 {instance.reply_to.hole} 的帖子 {instance.reply_to} 被回复了'
#             create_and_send_message(instance.user, message)
