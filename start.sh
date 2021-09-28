#!/bin/bash
pipenv shell
python manage.py migrate
python manage.py loaddata init_data
python start.py
gunicorn -c gconfig.py OpenTreeHole.asgi