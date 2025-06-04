FROM golang:1.24-alpine AS builder

RUN apk add --no-cache tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/main.go

FROM alpine:latest

RUN apk add --no-cache tzdata

ENV TZ=Europe/Moscow

COPY --from=builder /app/main /docker_app/main
COPY --from=builder /app/config.yaml /docker_app/config.yaml
COPY --from=builder /app/.env /docker_app/.env
COPY --from=builder /app/services/kafkaOrder/orders /docker_app/services/kafkaOrder/orders
COPY --from=builder /app/services/kafkaPayment/payments /docker_app/services/kafkaPayment/payments

WORKDIR /docker_app

EXPOSE 8080

CMD ["/docker_app/main"]
