# tgdeploy

Usually Telegram bots have alike requirements,
I somewhat embrace copy+paste method for this case.

Regenerate README.md: `go run ./cmd`.

## Dockerfile

```dockerfile
ARG VERSION_GOLANG="1.23"
ARG VERSION_ALPINE="3.21"

FROM golang:${VERSION_GOLANG}-alpine${VERSION_ALPINE} AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download
COPY . .

RUN go build -o main ./cmd/

FROM alpine:${VERSION_ALPINE}

RUN apk --no-cache add ca-certificates tzdata && \
    update-ca-certificates

COPY --from=builder /app/main .

ENTRYPOINT ["./main"]
```

## compose.yaml

```yaml
services:
  kittenbark_tg:
    build:
      context: .
      dockerfile: Dockerfile
    # container_name: kittenbark_tg
    # network_mode: host
    # env_file:
    #   - .env
    #   - ${HOME}/.env
    # volumes:
    #   - ${TG_LOCAL_API_STORAGE}:${TG_LOCAL_API_STORAGE}
    # restart: on-failure
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
    deploy:
      mode: replicated
      replicas: 3
      resources:
        limits:
          cpus: '0.50'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3
        window: 120s
```
