from django.conf import settings
from django.contrib.auth.models import User
from django.db import models
from django.db.models.signals import post_save
from django.dispatch import receiver
from rest_framework.authtoken.models import Token


# 自动在创建用户后创建其 Token 和用户资料数据
@receiver(post_save, sender=settings.AUTH_USER_MODEL)
def create_token_and_profile(sender, instance=None, created=False, **kwargs):
    if created:
        Token.objects.create(user=instance)
        Profile.objects.create(user=instance)


class Division(models.Model):
    name = models.CharField(max_length=32, unique=True)
    description = models.TextField(null=True)

    # def __str__(self):
    #     return self.name


class Tag(models.Model):
    name = models.CharField(max_length=8, unique=True)
    temperature = models.IntegerField(db_index=True, default=0, help_text="该标签下的主题帖数")

    def __str__(self):
        return self.name


# 主题帖
class Hole(models.Model):
    time_created = models.DateTimeField(auto_now_add=True)
    time_updated = models.DateTimeField(auto_now=True, db_index=True)
    count = models.IntegerField(db_index=True, default=0, help_text="楼层数")
    tags = models.ManyToManyField(Tag, blank=True)
    division = models.ForeignKey(Division, on_delete=models.CASCADE, help_text="分区")
    view = models.IntegerField(db_index=True, default=0, help_text="浏览量")
    deleted = models.BooleanField(default=False)
    mapping = models.JSONField(help_text='匿名到真实用户的对应')  # {user.id: anonymous_name}

    # key_floors 首条和末条回帖，动态生成

    def __str__(self):
        return "树洞#{}:{}".format(self.pk, self.floor_set.order_by("pk")[0].content[:50])


class Floor(models.Model):
    hole = models.ForeignKey(Hole, on_delete=models.CASCADE)
    content = models.TextField()
    anonyname = models.CharField(max_length=16)
    user = models.ForeignKey(User, models.CASCADE)
    reply_to = models.IntegerField(null=True)
    time_created = models.DateTimeField(auto_now_add=True)
    time_updated = models.DateTimeField(auto_now=True)
    like = models.IntegerField(default=0, db_index=True)
    deleted = models.BooleanField(default=False)
    delete_reason = models.TextField(null=True)
    folded = models.JSONField(null=True)

    def __str__(self):
        return "树洞#{}, 楼层#{}: {}".format(self.hole.pk, self.pk, self.content[:50])


# class Mapping(models.Model):
#     user = models.ForeignKey(User, on_delete=models.CASCADE)
#     anonyname = models.CharField(max_length=16)
#     hole = models.ForeignKey(Hole, on_delete=models.CASCADE)
#
#     def __str__(self):
#         return "#{}: {} -> {}".format(self.hole.pk, self.user.pk, self.anonyname)


class Report(models.Model):
    hole = models.ForeignKey(Hole, on_delete=models.CASCADE)
    floor = models.ForeignKey(Floor, on_delete=models.CASCADE)
    reason = models.TextField()
    time_created = models.DateTimeField(auto_now_add=True)
    dealed = models.BooleanField(default=False)
    dealed_by = models.ForeignKey(User, on_delete=models.CASCADE)

    def __str__(self):
        return "帖子#{}，{}".format(self.hole.pk, self.reason)


class Profile(models.Model):
    user = models.OneToOneField(User, on_delete=models.CASCADE)
    favorites = models.ManyToManyField(Hole, blank=True)
    permission = models.JSONField(null=True)

    # permission = {
    #     "admin": "$time",           # 管理员权限：到期时间
    #     "silent": {                 # 禁言     ：到期时间
    #         "$division_id": "$time" # 分区ID   ：到期时间
    #     }
    # }

    def __str__(self):
        return "用户数据#{}".format(self.user.pk)


class Message(models.Model):
    from_user = models.ForeignKey(
        User, related_name="message_from", on_delete=models.CASCADE, db_index=True
    )
    to_user = models.ForeignKey(
        User, related_name="message_to", on_delete=models.CASCADE, db_index=True
    )
    content = models.TextField()
    time_created = models.DateTimeField(auto_now_add=True)

    def __str__(self):
        return "{} -> {}: {}".format(self.from_user.pk, self.to_user.pk, self.content)
