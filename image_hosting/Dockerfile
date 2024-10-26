FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN apk add --no-cache --virtual .build-deps \
        ca-certificates \
        gcc \
        g++ &&  \
    go mod download

COPY . .

RUN go build -ldflags "-s -w" -o image_proxy

FROM alpine

RUN apk add --no-cache tzdata

ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/image_proxy /app/

EXPOSE 8000

ENTRYPOINT ["./image_proxy"]