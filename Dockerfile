# Этап 1: Сборка (build stage)
FROM golang:1.27-alpine AS builder

WORKDIR /app

# Копируем исходный код
COPY *.go .
COPY internal/ internal/
COPY go.mod .
COPY go.sum .

# Собираем бинарник
RUN go build -o app .

# Этап 2: Runtime (финальный образ)
FROM alpine:latest

WORKDIR /app

# Копируем только готовый бинарник из builder
COPY --from=builder /app/app .

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./app"]
