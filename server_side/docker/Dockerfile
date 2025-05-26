# Этап сборки
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

# Финальный минимальный образ
FROM alpine:latest

COPY --from=builder /app/config.yaml /docker_app/config.yaml
COPY --from=builder /app/main /docker_app/main

WORKDIR /docker_app

EXPOSE 8080

CMD ["/docker_app/main"]

