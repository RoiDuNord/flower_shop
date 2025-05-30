# Этап 1: Сборка
FROM golang:1.24-alpine AS builder

# Установка tzdata и зависимостей для сборки
RUN apk add --no-cache tzdata

WORKDIR /app

# Кэш зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код и собираем бинарник
COPY . .
RUN go build -o main ./cmd/main.go

# Этап 2: Финальный минимальный образ
FROM alpine:latest

# Установка только tzdata (для корректного времени в логах)
RUN apk add --no-cache tzdata

# Настройка временной зоны (по желанию, может быть изменена)
ENV TZ=Europe/Moscow

# Копируем нужные файлы из builder
COPY --from=builder /app/main /docker_app/main
COPY --from=builder /app/config.yaml /docker_app/config.yaml
COPY --from=builder /app/.env /docker_app/.env
COPY --from=builder /app/services/kafkaOrder/orders /docker_app/services/kafkaOrder/orders
COPY --from=builder /app/services/kafkaPayment/payments /docker_app/services/kafkaPayment/payments

WORKDIR /docker_app

# Открываем нужный порт
EXPOSE 8080

# Команда запуска
CMD ["/docker_app/main"]
