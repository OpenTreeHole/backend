FROM debian:buster

MAINTAINER jsclndnz@gmail.com

ENV HOLE_ENV=production REDIS_URL=redis://redis:6379 DEBIAN_FRONTEND=noninteractive

RUN apt update \
    && apt install -y lsb-release curl wget gnupg python3 python3-pip python3-dev libmagic1 \
    && curl -sLo mysql.deb https://dev.mysql.com/get/mysql-apt-config_0.8.19-1_all.deb \
    && DEBIAN_FRONTEND=noninteractive dpkg -i mysql.deb \
    && rm mysql.deb \
    && apt update \
    && apt install -y libmysqlclient-dev \
    && apt remove -y lsb-release curl wget gnupg \
    && apt autoremove -y \
    && apt clean \
    && pip3 install --no-cache-dir pipenv
    
WORKDIR /www/backend

COPY Pipfile /www/backend/

RUN pipenv install --skip-lock

COPY . /www/backend

EXPOSE 80

RUN chmod +x start.sh

ENTRYPOINT ["./start.sh"]
