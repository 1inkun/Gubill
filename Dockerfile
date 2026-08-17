# 多阶段构建：纯 Go SQLite 驱动，无需 CGO/gcc
FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/gubill ./cmd/main \
    && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/gubill-cli ./cmd/cli

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /out/gubill /app/gubill
COPY --from=builder /out/gubill-cli /app/gubill-cli
COPY cmd/main/config.yaml /app/config.yaml
COPY scripts/backup.sh /app/scripts/backup.sh
RUN chmod +x /app/scripts/backup.sh

ENV CONFIG_PATH=/app/config.yaml
VOLUME ["/app/data"]
EXPOSE 8080 8081

CMD ["./gubill"]
