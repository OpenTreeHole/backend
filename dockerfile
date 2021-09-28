FROM debian:buster

MAINTAINER jsclndnz@gmail.com

RUN apt update \
    && apt install -y --no-install-recommends python3 python3-pip libmagic1 \
    && apt autoremove -y \
	&& pip3 install --no-cache-dir pipenv

WORKDIR /www/backend

ENV HOLE_ENV=production

COPY Pipfile /www/backend/

RUN pipenv install --skip-lock

COPY . /www/backend

EXPOSE 80

ENTRYPOINT ["start.sh"]
