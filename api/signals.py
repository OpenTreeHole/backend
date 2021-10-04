from django.conf import settings
from django.db.models import F
from django.db.models.signals import post_save, m2m_changed, pre_save
from django.dispatch import receiver
from rest_framework.authtoken.models import Token

from api.models import Hole, Tag, Message, Floor
from api.utils import send_message_to_user, to_shadow_text


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


# 在数据库中创建一条消息并通过 websocket 发送给用户
def create_and_send_message(user, message):
    Message.objects.create(user=user, content=message)
    send_message_to_user(user, {'message': message})


# 收到回复后通知用户
@receiver(post_save, sender=Floor)
def notify_when_replied(sender, instance, created, **kwargs):
    if created:
        for floor in instance.mention.all():
            if 'reply' in floor.user.config['notify']:
                message = f'你在 {floor.hole} 的帖子 {floor} 被回复了'
                create_and_send_message(floor.user, message)

# 收藏的主题帖有新帖时通知用户
# @receiver(post_save, sender=Floor)
# def notify_when_favorites_updated(sender, instance, created, **kwargs):
#     if created and instance.hole:
#         if 'reply' in instance.mention.user.config['notify']:
#             message = f'你在 {instance.mention.hole} 的帖子 {instance.mention} 被回复了'
#             create_and_send_message(instance.user, message)
