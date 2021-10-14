FROM debian:buster

MAINTAINER jsclndnz@gmail.com

RUN apt update \
    && apt install -y lsb-release curl \
    && curl -sLo mysql.deb https://dev.mysql.com/get/mysql-apt-config_0.8.19-1_all.deb \
    && DEBIAN_FRONTEND=noninteractive dpkg -i mysql.deb \
    && rm mysql.deb \
    && apt update \
    && apt install -y libmysqlclient-dev \
    && apt install -y --no-install-recommends python3 python3-pip libmagic1 python3-dev \
    && pip3 install --no-cache-dir pipenv \
    && apt remove -y python3-pip curl lsb-release \
    && apt autoremove -y \
    
    
WORKDIR /www/backend

ENV HOLE_ENV=production REDIS_URL=redis://redis:6379

COPY Pipfile /www/backend/

RUN pipenv install --skip-lock

COPY . /www/backend

EXPOSE 80

RUN chmod +x start.sh

ENTRYPOINT ["./start.sh"]
