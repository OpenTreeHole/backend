from datetime import datetime, timezone

from rest_framework import permissions
from django.utils.dateparse import parse_datetime

MODIFY_METHODS = ('PUT', 'PATCH', 'DELETE')


def has_permission(user, category):
    """
    判断所给用户是否具有给定权限
    Args:
        user:       用户实例
        category:   'admin': 管理员权限
                    integer: 在分区id内发帖权限

    Returns: boolean

    """
    now = datetime.now(timezone.utc)

    if category == 'admin':
        expire_time = parse_datetime(user.profile.permission['admin'])
        return expire_time > now

    else:
        silent = user.profile.permission['silent']
        if not silent.get(category):  # 未设置禁言，返回 True
            return True
        else:
            expire_time = parse_datetime(silent.get(category))
            return expire_time < now


class OnlyAdminCanModify(permissions.BasePermission):
    """
    适用于主题帖
    """

    def has_permission(self, request, view):
        if request.method in MODIFY_METHODS:
            return has_permission(request.user, 'admin')
        else:
            return True


class OwnerOrAdminCanModify(permissions.BasePermission):
    """
    适用于回复帖或用户资料
    """

    def has_object_permission(self, request, view, obj):
        if request.method in MODIFY_METHODS:
            return obj.user == request.user or has_permission(request.user, 'admin')
        else:
            return True


class NotSilentOrAdminCanPost(permissions.BasePermission):
    """
    在给定分区内是否具有发帖权限，传入一个hole对象
    """

    def has_object_permission(self, request, view, division_id):
        if request.method == 'POST':
            return has_permission(request.user, division_id) or has_permission(request.user, 'admin')
        else:
            return True
