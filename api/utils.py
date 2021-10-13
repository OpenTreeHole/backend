import re
from functools import wraps
from io import StringIO

from django.core.cache import cache
from markdown import Markdown
from rest_framework.views import exception_handler


def encrypt_email(email_cleartext):
    return email_cleartext


def custom_exception_handler(exc, context):
    # Call REST framework's default exception handler first,
    # to get the standard error response.
    response = exception_handler(exc, context)

    # 默认错误消息字段改为“message”
    if response is not None and response.data.get('detail'):
        response.data['message'] = str(response.data['detail'])
        del (response.data['detail'])

    return response


def to_shadow_text(content):
    """
    Markdown to plain text
    """

    def unmark_element(element, stream=None):
        if stream is None:
            stream = StringIO()
        if element.text:
            stream.write(element.text)
        for sub in element:
            unmark_element(sub, stream)
        if element.tail:
            stream.write(element.tail)
        return stream.getvalue()

    # patching Markdown
    Markdown.output_formats["plain"] = unmark_element
    # noinspection PyTypeChecker
    md = Markdown(output_format="plain")
    md.stripTopLevelTags = False

    # 该方法会把 ![text](url) 中的 text 丢弃，因此需要手动替换
    content = re.sub(r'!\[(.+)]\(.+\)', r'\1', content)

    return md.convert(content)


def cache_function_call(key, timeout):
    def decorate(func):
        cache_key = f'cache-{key}'

        @wraps(func)
        def wrapper(*args, **kwargs):
            if cache.get(cache_key):
                print('hit cache', key)
                return cache.get(cache_key)
            else:
                result = func(*args, **kwargs)
                cache.set(cache_key, result, timeout)
                print('set cache', key)
                return result

        return wrapper

    return decorate
