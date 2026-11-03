# Этап 1: сборка
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Сначала зависимости, чтобы слой кешировался
COPY go.mod go.sum ./
RUN go mod download

# Потом весь исходный код
COPY . .

RUN CGO_ENABLED=0 go build -o app .

# Этап 2: финальный образ
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]
