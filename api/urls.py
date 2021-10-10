"""OpenTreeHole URL Configuration

The `urlpatterns` list routes URLs to views. For more information please see:
    https://docs.djangoproject.com/en/3.2/topics/http/urls/
Examples:
Function views
    1. Add an import:  from my_app import views
    2. Add a URL to urlpatterns:  path('', views.home, name='home')
Class-based views
    1. Add an import:  from other_app.views import Home
    2. Add a URL to urlpatterns:  path('', Home.as_view(), name='home')
Including another URLconf
    1. Import the include() function: from django.urls import include, path
    2. Add a URL to urlpatterns:  path('blog/', include('blog.urls'))
"""
# from django.contrib import admin
from django.urls import path

from api import consumers
from api.api import index, login, RegisterApi, HolesApi, FloorsApi, TagsApi, FavoritesApi, ReportsApi, ImagesApi, MessagesApi, UsersApi, DivisionsApi, logout, VerifyApi

websocket_urlpatterns = [
    path('ws/notification', consumers.NotificationConsumer.as_asgi()),
]

urlpatterns = [
    # path("admin/", admin.site.urls),
    path("", index),
    path("login", login),
    path('logout', logout),
    path("register", RegisterApi.as_view()),
    path("verify/<str:method>", VerifyApi.as_view()),
    path("holes", HolesApi.as_view()),
    path("holes/<int:hole_id>", HolesApi.as_view()),
    path('floors', FloorsApi.as_view()),
    path('floors/<int:floor_id>', FloorsApi.as_view()),
    path('tags', TagsApi.as_view()),
    path('tags/<int:tag_id>', TagsApi.as_view()),
    path('user/favorites', FavoritesApi.as_view()),
    path('reports', ReportsApi.as_view()),
    path('reports/<int:report_id>', ReportsApi.as_view()),
    path('images', ImagesApi.as_view()),
    path('messages', MessagesApi.as_view()),
    path('messages/<int:message_id>', MessagesApi.as_view()),
    path('users/<int:user_id>', UsersApi.as_view()),
    path('users', UsersApi.as_view()),
    path('divisions', DivisionsApi.as_view()),
    path('divisions/<int:division_id>', DivisionsApi.as_view()),
]
