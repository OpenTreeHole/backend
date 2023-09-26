FROM golang:1.21-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN apk add --no-cache --virtual .build-deps \
        ca-certificates \
        gcc \
        g++ &&  \
    go mod download

COPY . .

RUN go build -o treehole -ldflags "-s -w" ./cmd

FROM alpine

WORKDIR /app

COPY --from=builder /app/treehole /app/
COPY --from=builder /app/config/config_default.json /app/config/
VOLUME ["/app/data", "/app/config"]

ENV TZ=Asia/Shanghai MODE=production LOG_LEVEL=info PORT=8000

EXPOSE 8000

ENTRYPOINT ["./treehole"]