#!/bin/sh
python3 manage.py migrate
python3 start.py
if [ "$1" = "test" ]; then
  python3 manage.py test
else
  gunicorn -c gconfig.py OpenTreeHole.asgi
fi
