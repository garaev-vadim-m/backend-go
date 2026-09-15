# Docker Compose

## Простой сервер

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
```

- `build: .` — собирает образ из Dockerfile
- `ports` — маппит порты

```bash
docker-compose up       # поднять
docker-compose down     # остановить
```

## С базой данных

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=myapp
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

volumes:
  postgres_data:
```

- `postgres` — сервис БД (готовый образ, не собираем)
- `depends_on` — app ждет, пока postgres поднимется
- `environment` — переменные для PostgreSQL
- `volumes` — данные сохраняются между перезагрузками
