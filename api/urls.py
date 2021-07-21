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

from api.api import *
from rest_framework.authtoken import views

urlpatterns = [
    # path("admin/", admin.site.urls),
    path("", index),
    path("login", login),
    path("register", register),
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
]
