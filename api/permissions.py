from rest_framework import permissions
from rest_framework.permissions import SAFE_METHODS

MODIFY_METHODS = ('PUT', 'PATCH', 'DELETE')


class OnlyAdminCanModify(permissions.BasePermission):
    """
    适用于主题帖
    """

    def has_permission(self, request, view):
        if request.method in MODIFY_METHODS:
            return request.user.is_admin
        else:
            return True


class OwnerOrAdminCanModify(permissions.BasePermission):
    """
    适用于回复帖或用户资料
    """

    def has_object_permission(self, request, view, obj):
        if request.method in MODIFY_METHODS:
            return obj.user == request.user or request.user.is_admin
        else:
            return True


class NotSilentOrAdminCanPost(permissions.BasePermission):
    """
    在给定分区内是否具有发帖权限，传入一个 division_id
    """

    def has_object_permission(self, request, view, division_id):
        if request.method == 'POST':
            return not request.user.is_silenced(division_id) or request.user.is_admin
        else:
            return True


class AdminOrReadOnly(permissions.BasePermission):
    def has_permission(self, request, view):
        if request.method in SAFE_METHODS:
            return True
        else:
            return request.user.is_admin


class AdminOrPostOnly(permissions.BasePermission):
    def has_permission(self, request, view):
        if request.method in ('POST', 'OPTIONS'):
            return True
        else:
            return request.user.is_admin


class OwenerOrAdminCanSee(permissions.BasePermission):
    def has_object_permission(self, request, view, instance):
        if request.method == 'GET':
            return instance.user == request.user or request.user.is_admin
        else:
            return True
