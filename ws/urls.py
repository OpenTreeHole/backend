from django.urls import path

from ws.notification import NotificationConsumer

urlpatterns = [
    path('ws/notification', NotificationConsumer.as_asgi()),
]
