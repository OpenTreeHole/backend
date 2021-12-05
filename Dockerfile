FROM python:3.9 as builder

ENV PIPENV_VENV_IN_PROJECT="enabled" DEBIAN_FRONTEND=noninteractive

WORKDIR /www/backend

RUN apt update \
    && apt install -y default-libmysqlclient-dev python3-dev libmagic1 \
    && pip3 install --no-cache-dir pipenv

COPY Pipfile /www/backend/

RUN pipenv install --skip-lock

FROM python:3.9-slim

MAINTAINER jsclndnz@gmail.com

RUN apt update \
    && apt install -y --no-install-recommends default-libmysqlclient-dev libmagic1 \
    && apt autoremove && apt clean && rm -rf /var/lib/apt/lists/*

WORKDIR /www/backend

COPY --from=builder /www/backend/.venv /www/backend/.venv

COPY . /www/backend

ENV HOLE_ENV=production \
    REDIS_URL=redis://redis:6379 \
    PATH="/www/backend/.venv/bin:$PATH"

EXPOSE 80

RUN chmod +x start.sh \
# APNS 证书问题临时解决方案，安全性低
    && echo "CipherString=DEFAULT@SECLEVEL=1" >> /etc/ssl/openssl.cnf

ENTRYPOINT ["./start.sh"]
