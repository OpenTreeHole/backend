#!/bin/sh
pipenv run python manage.py migrate
pipenv run python manage.py loaddata init_data
pipenv run python start.py
pipenv run gunicorn -c gconfig.py OpenTreeHole.asgi
