"""
ASGI config for OpenTreeHole project.

It exposes the ASGI callable as a module-level variable named ``application``.

For more information on this file, see
https://docs.djangoproject.com/en/3.2/howto/deployment/asgi/
"""

import os

from django.core.asgi import get_asgi_application

django_asgi_app = get_asgi_application()  # 必须写在所有导入前面

from channels.routing import ProtocolTypeRouter, URLRouter
from api.middleware import TokenAuthMiddleware

import api.urls

os.environ.setdefault("DJANGO_SETTINGS_MODULE", "OpenTreeHole.settings")

application = ProtocolTypeRouter({
    "http": django_asgi_app,
    "websocket": TokenAuthMiddleware(
        URLRouter(
            api.urls.websocket_urlpatterns
        )
    ),
})
