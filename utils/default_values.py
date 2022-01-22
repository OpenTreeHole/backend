"""
数据库默认值生成函数
"""
from datetime import datetime, timedelta

from django.conf import settings

from utils.constants import NotifyConfig


def default_active_user_date():
    return datetime.now(settings.TIMEZONE) - timedelta(days=1)


def default_permission():
    """
    silent 字典
        index：分区id （string） django的JSONField会将字典的int索引转换成str
        value：禁言解除时间
    """
    return {
        'admin': '1970-01-01T00:00:00+00:00',  # 管理员权限：到期时间
        'silent': {},  # 禁言
        'offense_count': 0
    }


def default_config():
    """
    show_folded: 对折叠内容的处理
        fold: 折叠
        hide: 隐藏
        show: 展示

    notify: 在以下场景时通知
        NotifyConfig.floor_mentioned:       帖子被提及时
        NotifyConfig.favored_hole_replied:  收藏的主题帖有新帖时
        NotifyConfig.reported:              被举报时通知管理员
        NotifyConfig.punished:              被处罚时
    另外，当用户权限发生变化或所发帖被修改时也会收到通知
    """
    return {
        'show_folded': 'fold',
        'notify': [NotifyConfig.floor_mentioned, NotifyConfig.favored_hole_replied,
                   NotifyConfig.punished]
    }
