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
from rest_framework_simplejwt.views import (
    TokenObtainPairView,
    TokenRefreshView,
)

from api.api import index, login, RegisterApi, verify, HolesApi, FloorsApi, TagsApi, FavoritesApi, ReportsApi, ImagesApi

urlpatterns = [
    # path("admin/", admin.site.urls),
    path("", index),
    path('token', TokenObtainPairView.as_view(), name='token_obtain_pair'),
    path('token/refresh', TokenRefreshView.as_view(), name='token_refresh'),
    path("login", login),
    path("register", RegisterApi.as_view()),
    path("verify/<str:method>", verify),
    path("holes", HolesApi.as_view()),
    path("holes/<int:hole_id>", HolesApi.as_view()),
    path('floors', FloorsApi.as_view()),
    path('floors/<int:floor_id>', FloorsApi.as_view()),
    path('tags', TagsApi.as_view()),
    path('tags/<int:tag_id>', TagsApi.as_view()),
    path('user/favorites', FavoritesApi.as_view()),
    path('reports', ReportsApi.as_view()),
    path('reports/<int:report_id>', ReportsApi.as_view()),
    path('images', ImagesApi.as_view())
]
