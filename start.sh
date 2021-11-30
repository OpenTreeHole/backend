#!/bin/sh
poetry run python3 manage.py migrate
poetry run python3 start.py
poetry run gunicorn -c gconfig.py OpenTreeHole.asgi
