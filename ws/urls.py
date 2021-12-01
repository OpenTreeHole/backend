from django.urls import path

from ws import notification

urlpatterns = [
    path('ws/notification', notification.NotificationConsumer.as_asgi()),
]
