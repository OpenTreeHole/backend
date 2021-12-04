"""
业务相关
"""

import re
from io import StringIO

from django.http import Http404
from markdown import Markdown


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


def exists_or_404(klass, *args, **kwargs):
    if hasattr(klass, '_default_manager'):
        # noinspection PyProtectedMember
        if not klass._default_manager.filter(*args, **kwargs).exists():
            raise Http404(f'{klass} 对象不存在！')
    else:
        klass__name = klass.__name__ if isinstance(klass, type) else klass.__class__.__name__
        raise ValueError(
            "First argument to get_object_or_404() must be a Model, Manager, "
            "or QuerySet, not '%s'." % klass__name
        )