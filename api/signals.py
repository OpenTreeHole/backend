from django.conf import settings
from django.db.models.signals import post_save
from django.dispatch import receiver
from rest_framework.authtoken.models import Token

from api.models import Profile


# 自动在创建用户后创建其 Token 和用户资料数据
@receiver(post_save, sender=settings.AUTH_USER_MODEL)
def create_token_and_profile(sender, instance=None, created=False, **kwargs):
    print('hi')
    if created:
        print('hi')
        Token.objects.create(user=instance)
        Profile.objects.create(user=instance)
