# hello-go-backend

Небольшой Go-бэкенд с PostgreSQL и Docker.

## Быстрый старт

```bash
docker compose up --build
```

## Веб-страницы

- `index.html` — простая витрина интернет-магазина с карточками товаров (данные из `/products`).
- `admin_panel.html` — самодостаточная страница (ванильный JS) для работы с API: CRUD для `users`, `status`, `rules`, `user_status`, `products`, проверка `/health`.
- `auth.html` — форма входа по email и паролю (`/user/login`).

После успешной авторизации сессия сохраняется в `localStorage`. Пользователь с ролью `admin` перенаправляется в `admin_panel.html`, пользователь с другой ролью остаётся на странице авторизации. Операции создания/изменения/удаления в панели доступны только пользователю с ролью `admin` (проверка только на стороне панели, сервер права пока не проверяет). Панель проверяет длину пароля (8–72 байта) и лимиты полей из схемы БД, а при `401`/`403` очищает сессию и возвращает на `auth.html`.

Страницы обращаются к API по `API base URL` (по умолчанию `http://localhost:8080`; значение хранится в `localStorage`). Сервер отдаёт CORS-заголовки (`corsMiddleware` в `main.go`) и отвечает `204` на preflight `OPTIONS`.

## Локальная разработка (без Docker)

1. Поднимите PostgreSQL и создайте БД `myapp`
2. Экспортируйте переменные окружения:

```bash
export DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=myapp
```

Подключение к базе в контейнере:
```bash
docker compose exec postgres psql -U postgres -d myapp
```

Проверить таблицы:
```bash
docker compose exec postgres psql -U postgres -d myapp -c "\dt"
```

3. Запуск:

```bash
go run .
```

## Стек

- Go 1.22
- PostgreSQL 15
- Docker / Docker Compose
