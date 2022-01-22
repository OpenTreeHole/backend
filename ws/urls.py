from django.urls import path

from ws.images import ImageConsumer
from ws.notification import NotificationConsumer

urlpatterns = [
    path('ws/notification', NotificationConsumer.as_asgi()),
    path('ws/images', ImageConsumer.as_asgi()),
]
