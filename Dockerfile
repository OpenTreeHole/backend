FROM golang:1.22-alpine as builder

ARG SERVICE_NAME

WORKDIR /app

COPY go.mod go.sum ./

RUN apk add --no-cache ca-certificates tzdata && \
    go mod download

COPY . .

RUN go build -ldflags "-s -w" -tags netgo -o backend ./$SERVICE_NAME/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /app/backend /app/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

VOLUME ["/app/data"]

ENV TZ=Asia/Shanghai
ENV MODE=prod
ENV LOG_LEVEL=info

EXPOSE 8000

ENTRYPOINT ["./backend"]