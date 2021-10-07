#!/bin/sh
pipenv run python manage.py migrate
pipenv run python start.py
pipenv run gunicorn -c gconfig.py OpenTreeHole.asgi
