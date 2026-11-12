# Docker Compose

## Простой сервер

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
```

- `build: .` — собирает образ из Dockerfile
- `ports` — маппит порты

```bash
docker compose up --build   # собрать и поднять
docker compose up -d        # поднять в фоне
docker compose down         # остановить и удалить
```

`docker compose` (v2) — актуальная команда. Старый вариант `docker-compose` с дефисом тоже работает, если установлен отдельно.

## С базой данных

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      DB_HOST: postgres
      DB_PORT: "5432"
      DB_USER: postgres
      DB_PASSWORD: postgres
      DB_NAME: myapp
    depends_on:
      - postgres

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
```

- `postgres` — сервис БД (готовый образ, не собираем)
- `depends_on` — app стартует после postgres
- `environment` — переменные окружения; `DB_HOST` указывает на имя сервиса `postgres`
- `volumes` — данные сохраняются между перезапусками
- `healthcheck` — проверка, что PostgreSQL готов принимать подключения

## healthcheck

```yaml
healthcheck:
  test: ["CMD-SHELL", "pg_isready -U postgres"]
  interval: 10s
  timeout: 5s
  retries: 5
```

- `test` — команда проверки (выполняется внутри контейнера postgres)
- `interval` — как часто проверять
- `timeout` — сколько ждать ответа
- `retries` — сколько неудач подряд считать «контейнер нездоров»

Приложение всё равно не должно падать, если БД ещё поднимается: `InitDB` в `internal/database/database.go` делает 5 попыток подключения с паузой 2 секунды.

## Ключи compose

| Ключ          | Значение                              |
|---------------|---------------------------------------|
| `build`       | собрать образ из Dockerfile           |
| `image`       | использовать готовый образ            |
| `ports`       | пробросить порт `хост:контейнер`      |
| `environment` | переменные окружения сервиса          |
| `depends_on`  | порядок запуска сервисов              |
| `volumes`     | постоянное хранилище данных           |
| `healthcheck` | проверка готовности сервиса           |

```bash
docker compose up --build   # поднять всё
docker compose logs -f app  # логи приложения
docker compose down         # остановить и удалить контейнеры
```
