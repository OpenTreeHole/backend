from django.conf import settings
from django.contrib.auth.models import User
from django.db import models
from django.db.models import F
from django.db.models.signals import post_save, m2m_changed
from django.dispatch import receiver
from rest_framework.authtoken.models import Token


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
        return f'树洞#{self.pk}: {self.floor_set.order_by("pk")[0].content[:50]}'


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
    user = models.ForeignKey(User, models.CASCADE)
    reply_to = models.IntegerField(null=True)
    time_created = models.DateTimeField(auto_now_add=True)
    time_updated = models.DateTimeField(auto_now=True)
    like = models.IntegerField(default=0, db_index=True)  # 赞同数
    like_data = models.JSONField(default=list)  # 点赞记录，主键列表
    deleted = models.BooleanField(default=False)  # 仅作为前端是否显示删除按钮的依据
    history = models.JSONField(default=list)  # 修改记录，字典列表
    fold = models.JSONField(default=list)  # 折叠原因，字符串列表（原因由前端提供）

    def __str__(self):
        return f"树洞#{self.hole.pk}, 楼层#{self.pk}: {self.content[:50]}"


class Report(models.Model):
    hole = models.ForeignKey(Hole, on_delete=models.CASCADE)
    floor = models.ForeignKey(Floor, on_delete=models.CASCADE)
    reason = models.TextField()
    time_created = models.DateTimeField(auto_now_add=True)
    time_updated = models.DateTimeField(auto_now=True)
    dealed = models.BooleanField(default=False, db_index=True)
    dealed_by = models.ForeignKey(User, on_delete=models.CASCADE, null=True)

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


class Profile(models.Model):
    user = models.OneToOneField(User, on_delete=models.CASCADE)
    nickname = models.CharField(max_length=32, blank=True)
    favorites = models.ManyToManyField(Hole, blank=True)
    permission = models.JSONField(default=default_permission)
    config = models.JSONField(default=default_config)

    def __str__(self):
        return f"用户数据#{self.user.pk}"


class Message(models.Model):
    user = models.ForeignKey(
        User, related_name="message_to", on_delete=models.CASCADE, db_index=True
    )
    content = models.TextField()
    has_read = models.BooleanField(default=False)
    time_created = models.DateTimeField(auto_now_add=True)

    def __str__(self):
        return f"-> {self.user.pk}: {self.content}"


# 自动在创建用户后创建其 Token 和用户资料数据
@receiver(post_save, sender=settings.AUTH_USER_MODEL)
def create_token_and_profile(sender, instance=None, created=False, **kwargs):
    if created:
        Token.objects.create(user=instance)
        Profile.objects.create(user=instance)


# 自动修改 tag 的热度
@receiver(m2m_changed, sender=Hole.tags.through)
def modify_tag_temperature(sender, reverse, action, pk_set, **kwargs):
    if reverse is False and action == 'post_add':
        Tag.objects.filter(pk__in=pk_set).update(temperature=F('temperature') + 1)
    elif reverse is False and action == 'post_remove':
        Tag.objects.filter(pk__in=pk_set).update(temperature=F('temperature') - 1)
