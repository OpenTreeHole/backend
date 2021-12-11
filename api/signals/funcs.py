from datetime import datetime, timezone

from django.conf import settings
from django.contrib.auth import get_user_model
from django.core.cache import cache
from django.db.models import F
from django.db.models.signals import post_save, m2m_changed, pre_save
from django.dispatch import receiver
from rest_framework.authtoken.models import Token

from api.models import Hole, Tag, Floor, Report
from api.serializers import FloorSerializer, ReportSerializer
from api.signals import modified_by_admin, mention_to, new_penalty
from utils.apis import to_shadow_text
from utils.notification import send_notifications


# 在创建用户后创建其 Token
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


# 添加 / 修改帖子后
@receiver(post_save, sender=Floor)
def after_adding_a_floor(sender, instance, created, **kwargs):
    # 添加帖子后增加 reply 数
    if created:
        Hole.objects.filter(id=instance.hole_id).update(reply=F('reply') + 1)
    # 修改帖子后清除缓存
    cache_key = f'_cached-{instance.hole}'
    cache.delete(cache_key)


# 帖子被提及后通知用户
@receiver(mention_to, sender=Floor)
def notify_when_mentioned(sender, instance, mentioned, **kwargs):
    """
    Args:
        instance: 提及的帖子
        sender: Floor
        mentioned: [<Floor: 被提及的帖子>]
        **kwargs:
    """
    for floor in mentioned:
        if 'mention' in floor.user.config['notify']:
            message = f'你在树洞#{floor.hole_id}的帖子##{floor.id}被引用了'
            data = FloorSerializer(instance, context={"user": floor.user}).data
            send_notifications.delay(floor.user_id, message, data, 'mention')


# 收藏的主题帖有新帖时通知用户
@receiver(post_save, sender=Floor)
def notify_when_favorites_updated(sender, instance, created, **kwargs):
    if created:
        for user in instance.hole.favored_by.filter(config__notify__icontains='favorite'):
            message = f'你收藏的树洞#{instance.hole_id}有新回复'
            data = FloorSerializer(instance, context={"user": user}).data
            send_notifications.delay(user.id, message, data, 'favorite')


# 被举报时通知用户和管理员
@receiver(post_save, sender=Report)
def notify_when_reported(sender, instance, created, **kwargs):
    floor = instance.floor
    if created:
        data = ReportSerializer(instance).data
        if 'report' in floor.user.config['notify']:
            message = f'你的帖子#{instance.hole}(##{instance.floor})被举报了'
            send_notifications.delay(floor.user_id, message, data, 'report')
        # 通知管理员
        queryset = get_user_model().objects.filter(permission__admin__gt=datetime.now(timezone.utc).isoformat()).values_list('id', flat=True)
        for admin_id in list(queryset):
            message = f'{floor.user}的树洞#{instance.hole}(##{instance.floor})被举报了'
            send_notifications.delay(admin_id, message, data, 'report')


# 用户权限发生变化时发送通知
@receiver(post_save, sender=settings.AUTH_USER_MODEL)
def notify_when_permission_changed(sender, instance, **kwargs):
    update_fields = kwargs.get('update_fields') if kwargs.get('update_fields') else []
    if 'permission' in update_fields:
        message = '你的权限被更改了'
        data = instance.permission
        send_notifications.delay(instance.id, message, data, 'permission')


# 用户帖子被修改后发出通知
@receiver(modified_by_admin, sender=Floor)
def notify_when_floor_modified_by_admin(sender, instance, **kwargs):
    message = f'你的帖子##{instance}被修改了'
    data = FloorSerializer(instance, context={"user": instance.user}).data
    send_notifications.delay(instance.user_id, message, data, 'modify')


# 用户被处罚后发送通知
@receiver(new_penalty, sender=Floor)
def notify_when_floor_modified_by_admin(sender, instance, penalty, **kwargs):
    message = f'你因为帖子##{instance}违规而被处罚'
    data = {
        "level": penalty[0],
        "date": penalty[1],
        "division_id": penalty[2],
    }
    send_notifications.delay(instance.user_id, message, data, 'penalty')
