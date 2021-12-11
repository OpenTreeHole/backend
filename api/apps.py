from django.apps import AppConfig


class ApiConfig(AppConfig):
    name = "api"

    # 需要在此处导入信号模块
    def ready(self):
        # noinspection PyUnresolvedReferences
        import api.signals.funcs
