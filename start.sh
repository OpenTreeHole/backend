#!/bin/sh
python3 manage.py migrate
python3 start.py
gunicorn -c gconfig.py OpenTreeHole.asgi
