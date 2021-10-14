FROM debian:buster

MAINTAINER jsclndnz@gmail.com

ENV HOLE_ENV=production REDIS_URL=redis://redis:6379 DEBIAN_FRONTEND=noninteractive

RUN apt update \
    && apt install -y lsb-release curl wget \
    && curl -sLo mysql.deb https://dev.mysql.com/get/mysql-apt-config_0.8.19-1_all.deb \
    && dpkg -i mysql.deb \
    && rm mysql.deb \
    && apt update \
    && apt install -y libmysqlclient-dev \
    && apt remove -y curl wget lsb-release \
    && apt autoremove -y \
    && apt clean
    
RUN apt install -y --no-install-recommends python3 python3-pip libmagic1 python3-dev \
    && pip3 install --no-cache-dir pipenv \
    && apt remove -y python3-pip \
    && apt autoremove -y \
    && apt clean
    
    
WORKDIR /www/backend


COPY Pipfile /www/backend/

RUN pipenv install --skip-lock

COPY . /www/backend

EXPOSE 80

RUN chmod +x start.sh

ENTRYPOINT ["./start.sh"]
