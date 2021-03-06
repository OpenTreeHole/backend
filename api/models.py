from datetime import datetime

from django.conf import settings
from django.contrib.auth.base_user import AbstractBaseUser, BaseUserManager
from django.db import models
from django.utils.dateparse import parse_datetime
from rest_framework.authtoken.models import Token

from utils.auth import encrypt_email, many_hashes
from utils.default_values import default_active_user_date, default_permission, default_config


class Division(models.Model):
    name = models.CharField(max_length=32, unique=True, default='树洞')
    description = models.TextField(null=True)
    pinned = models.JSONField(default=list)

    def __str__(self):
        return self.name

    class Meta:
        db_table = 'division'


class Tag(models.Model):
    name = models.CharField(max_length=settings.MAX_TAG_LENGTH, unique=True)
    temperature = models.IntegerField(db_index=True, default=0, help_text="该标签下的主题帖数")

    def __str__(self):
        return self.name

    class Meta:
        db_table = 'tag'


# 主题帖
class Hole(models.Model):
    time_created = models.DateTimeField(auto_now_add=True)
    time_updated = models.DateTimeField(auto_now=True, db_index=True)
    tags = models.ManyToManyField(Tag, blank=True)
    division = models.ForeignKey(Division, on_delete=models.CASCADE, help_text="分区")
    view = models.IntegerField(db_index=True, default=0, help_text="浏览量")
    reply = models.IntegerField(db_index=True, default=-1,
                                help_text="楼层数")  # 如果只有首条帖子的话认为回复数为零
    mapping = models.JSONField(default=dict,
                               help_text='匿名到真实用户的对应')  # {user.id: anonymous_name}
    hidden = models.BooleanField(default=False, help_text='帖子是否隐藏')

    def __str__(self):
        return f'树洞#{self.pk}'

    class Meta:
        db_table = 'hole'


class Floor(models.Model):
    """
    history:
        'content': floor.content,                               # 原内容
        'altered_by': request.user.pk,                          # 修改者 id
        'altered_time': datetime.now(settings.TIMEZONE).isoformat()  # 修改时间
    """
    hole = models.ForeignKey(Hole, on_delete=models.CASCADE)
    content = models.TextField()
    shadow_text = models.TextField(default='')  # 去除markdown关键字的文本，方便搜索
    anonyname = models.CharField(max_length=16)
    user = models.ForeignKey(settings.AUTH_USER_MODEL, models.CASCADE)
    mention = models.ManyToManyField('self', blank=True, symmetrical=False,
                                     related_name='mentioned_by')
    time_created = models.DateTimeField(auto_now_add=True)
    time_updated = models.DateTimeField(auto_now=True)
    like = models.IntegerField(default=0, db_index=True)  # 赞同数
    like_data = models.JSONField(default=list)  # 点赞记录，主键列表
    deleted = models.BooleanField(default=False)  # 仅作为前端是否显示删除按钮的依据
    history = models.JSONField(default=list)  # 修改记录，字典列表
    fold = models.JSONField(default=list)  # 折叠原因，字符串列表（原因由前端提供）
    special_tag = models.CharField(max_length=16, default='')  # 额外字段
    storey = models.IntegerField(default=0)  # 楼层数

    def __str__(self):
        return f"{self.content[:50]}"

    class Meta:
        db_table = 'floor'


class Report(models.Model):
    hole = models.ForeignKey(Hole, on_delete=models.CASCADE)
    floor = models.ForeignKey(Floor, on_delete=models.CASCADE)
    reason = models.TextField()
    time_created = models.DateTimeField(auto_now_add=True)
    time_updated = models.DateTimeField(auto_now=True)
    dealed = models.BooleanField(default=False, db_index=True)
    dealed_by = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.CASCADE,
                                  null=True)

    def __str__(self):
        return f"{self.hole}, 帖子{self.floor}\n理由: {self.reason}"

    class Meta:
        db_table = 'report'


class UserManager(BaseUserManager):
    def create_user(self, email, password=None, **extra_fields):
        """
        Args:
            email: 明文
            password: 明文
        Returns:
            user
        """
        if not email:
            raise ValueError('邮箱必须提供')
        email = self.normalize_email(email)
        user = self.model(
            email=encrypt_email(email),
            identifier=many_hashes(email),
            **extra_fields
        )
        user.set_password(password)
        user.save()
        return user

    def create_superuser(self, email, password=None, **extra_fields):
        user = self.create_user(email, password, **extra_fields)
        user.permission['admin'] = settings.VERY_LONG_TIME
        user.save()
        return user


class User(AbstractBaseUser):
    email = models.CharField(max_length=1000)  # RSA encrypted email
    identifier = models.CharField(max_length=128, unique=True)  # sha512 of email
    joined_time = models.DateTimeField(auto_now_add=True)
    nickname = models.CharField(max_length=32, blank=True)
    favorites = models.ManyToManyField(Hole, related_name='favored_by', blank=True)
    permission = models.JSONField(default=default_permission)
    config = models.JSONField(default=default_config)

    objects = UserManager()

    USERNAME_FIELD = 'identifier'

    @property
    def is_admin(self):
        now = datetime.now(settings.TIMEZONE)
        expire_time = parse_datetime(self.permission['admin'])
        return expire_time > now

    def is_silenced(self, division_id):
        now = datetime.now(settings.TIMEZONE)
        silent = self.permission['silent']
        division = str(division_id)  # JSON 序列化会将字典的 int 索引转换成 str
        if not silent.get(division):  # 未设置禁言，返回 False
            return False
        else:
            expire_time = parse_datetime(silent.get(division))
            return expire_time > now

    def __str__(self):
        return f"用户#{self.pk}"

    class Meta:
        db_table = 'user'


class Message(models.Model):
    user = models.ForeignKey(
        settings.AUTH_USER_MODEL, related_name="message_to", on_delete=models.CASCADE,
        db_index=True
    )
    message = models.TextField(default='')
    code = models.CharField(max_length=30, default='')
    data = models.JSONField(default=dict)
    has_read = models.BooleanField(default=False)
    time_created = models.DateTimeField(auto_now_add=True)

    def __str__(self):
        return self.message

    class Meta:
        db_table = 'message'


class PushToken(models.Model):
    user = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.CASCADE,
                             related_name='push_tokens')
    service = models.CharField(max_length=16, db_index=True)  # apns or mipush
    device_id = models.CharField(max_length=128, unique=True)
    token = models.CharField(max_length=128)

    class Meta:
        db_table = 'push_token'


class OldUserFavorites(models.Model):
    uid = models.CharField(max_length=11)
    favorites = models.JSONField()


class ActiveUser(models.Model):
    date = models.DateField(default=default_active_user_date, unique=True)
    dau = models.IntegerField(default=0)  # 日活
    mau = models.IntegerField(default=0)  # 月活

    class Meta:
        db_table = 'active_user'
