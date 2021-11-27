from datetime import datetime, timezone

from django.conf import settings
from django.contrib.auth import get_user_model
from django.core.cache import cache
from django.db.models import F
from django.db.models.signals import post_save, m2m_changed, pre_save
from django.dispatch import receiver, Signal
from rest_framework.authtoken.models import Token

from api.models import Hole, Tag, Floor, Report
from api.notification import MessageSender
from api.serializers import FloorSerializer, ReportSerializer
from utils.apis import to_shadow_text

modified_by_admin = Signal(providing_args=['instance'])


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


# 创建 shadow_text
@receiver(pre_save, sender=Floor)
def create_shadow_text(sender, instance, **kwargs):
    instance.shadow_text = to_shadow_text(instance.content)


# 添加回复帖后执行任务
@receiver(post_save, sender=Floor)
def after_adding_a_floor(sender, instance, created, **kwargs):
    if created:
        # Hole.objects.filter(id=instance.hole_id).update(reply=F('reply') + 1)  # 不知道为什么它运行无效
        instance.hole.reply += 1
        instance.hole.save()
        cache_key = f'cache-{instance.hole}'
        deleted = cache.delete(cache_key)


# 帖子被提及后通知用户
@receiver(post_save, sender=Floor)
def notify_when_mentioned(sender, instance, created, **kwargs):
    if created:
        message_sender = MessageSender()
        for floor in instance.mention.all():
            if 'mention' in floor.user.config['notify']:
                message = f'你在{floor.hole} 的帖子#{floor.id} 被提及了'
                data = FloorSerializer(floor, context={"user": floor.user}).data
                message_sender.create_and_queue_or_send_message(floor.user, message, data, 'mention')
        message_sender.commit()


# 收藏的主题帖有新帖时通知用户
@receiver(post_save, sender=Floor)
def notify_when_favorites_updated(sender, instance, created, **kwargs):
    if created:
        message_sender = MessageSender()
        for user in instance.hole.favored_by.all():
            if 'favorite' in user.config['notify']:
                message = f'你收藏的{instance.hole} 被回复了{instance.id}'
                data = FloorSerializer(instance, context={"user": user}).data
                message_sender.create_and_queue_or_send_message(user, message, data, 'favorite')
        message_sender.commit()


# 被举报时通知用户和管理员
@receiver(post_save, sender=Report)
def notify_when_reported(sender, instance, created, **kwargs):
    floor = instance.floor
    if created:
        message_sender = MessageSender()
        data = ReportSerializer(instance).data
        if 'report' in floor.user.config['notify']:
            message = f'你被举报了{instance.id}'
            message_sender.create_and_queue_or_send_message(floor.user, message, data, 'report')
        for admin in get_user_model().objects.filter(permission__admin__gt=datetime.now(timezone.utc).isoformat()):
            message = f'{floor.user}被举报了{instance.id}'
            message_sender.create_and_queue_or_send_message(admin, message, data, 'report')
        message_sender.commit()


# 用户权限发生变化时发送通知
@receiver(post_save, sender=settings.AUTH_USER_MODEL)
def notify_when_permission_changed(sender, instance, **kwargs):
    update_fields = kwargs.get('update_fields') if kwargs.get('update_fields') else []
    if 'permission' in update_fields:
        message = '你的权限被更改了'
        data = instance.permission
        MessageSender(instance, message, data, 'permission').commit()


# 用户帖子被修改后发出通知
@receiver(modified_by_admin, sender=Floor)
def notify_when_floor_modified_by_admin(sender, instance, **kwargs):
    message = '你的帖子被修改了'
    data = FloorSerializer(instance, context={"user": instance.user}).data
    MessageSender(instance.user, message, data, 'modify').commit()
