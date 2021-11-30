#!/bin/sh
python3 manage.py migrate
python3 start.py
pipenv run gunicorn -c gconfig.py OpenTreeHole.asgi
