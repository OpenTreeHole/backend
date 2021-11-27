from functools import wraps

from django.core.cache import cache


def cache_function_call(key, timeout):
    def decorate(func):
        cache_key = f'cache-{key}'

        @wraps(func)
        def wrapper(*args, **kwargs):
            if cache.get(cache_key):
                return cache.get(cache_key)
            else:
                result = func(*args, **kwargs)
                cache.set(cache_key, result, timeout)
                return result

        return wrapper

    return decorate
