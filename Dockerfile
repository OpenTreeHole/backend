FROM python:3.9 as builder

ENV PIPENV_VENV_IN_PROJECT="enabled"

WORKDIR /www/backend

RUN apt update \
    && apt install -y default-libmysqlclient-dev python3-dev libmagic1 python3-distutils \
    && pip3 install pipenv

RUN pipenv install \
    && cp -r /usr/local/lib/python3.9/site-packages/_distutils_hack .venv/lib/python3.9/site-packages

COPY Pipfile /www/backend/

RUN pipenv install --dev --skip-lock

FROM python:3.9-slim

MAINTAINER jsclndnz@gmail.com

RUN apt update \
    && apt install -y --no-install-recommends default-libmysqlclient-dev libmagic1 \
    && apt autoremove && apt clean && rm -rf /var/lib/apt/lists/*

WORKDIR /www/backend

COPY --from=builder /www/backend/.venv /www/backend/.venv

COPY . /www/backend

ENV PATH="/www/backend/.venv/bin:$PATH"

EXPOSE 80

RUN chmod +x start.sh \
# APNS 证书问题临时解决方案，安全性低
    && echo "CipherString=DEFAULT@SECLEVEL=1" >> /etc/ssl/openssl.cnf

ENTRYPOINT ["./start.sh"]
