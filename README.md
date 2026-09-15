# hello-go-backend

Минимальный Go-бэкенд с PostgreSQL в Docker.

## Быстрый старт

```bash
docker compose up --build
```

## Эндпоинты

| Маршрут     | Описание            |
|-------------|---------------------|
| `GET /`     | Hello World         |
| `GET /health` | Healthcheck (ping DB) |

## Локальная разработка (без Docker)

1. Запустить PostgreSQL и создать БД `myapp`
2. Экспортировать переменные окружения:

```bash
export DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=myapp
```

3. Запустить:

```bash
go run .
```

## Стек

- Go 1.27
- PostgreSQL 15
- Docker / Docker Compose
