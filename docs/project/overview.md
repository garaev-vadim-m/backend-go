# Проект hello-go-backend

## Что это

REST API на Go с PostgreSQL и тремя статическими страницами без сборки:
витрина товаров (`index.html`), вход/регистрация (`auth.html`) и админ-панель
(`admin_panel.html`). Страницы встроены в бинарник через `go:embed web` и
раздаются тем же сервером, что и API. API реализует CRUD для пяти сущностей
и аутентификацию на серверных сессиях.

- API и страницы: `http://localhost:8080`
- Статика: `web/index.html`, `web/auth.html`, `web/admin_panel.html`
- Документация: `docs/`

> Сервер пока не проверяет токен на CRUD-роутах: права admin реализованы только
> в админ-панели (скрытие кнопок изменения). Планы по middleware прав — в
> `docs/project/tasks.md`.

## Стек

| Компонент   | Версия / технология                |
|-------------|------------------------------------|
| Язык        | Go 1.22                            |
| БД          | PostgreSQL 15                      |
| Драйвер БД  | `github.com/lib/pq`                |
| Миграции    | `github.com/golang-migrate/migrate/v4` (SQL в `embed.FS`) |
| Пароли      | `golang.org/x/crypto/bcrypt`       |
| Контейнеры  | Docker, Docker Compose             |
| Фронтенд    | HTML + vanilla JS (без сборки)     |

## Архитектура

```
main.go
   ├─ database.InitDB / RunMigrations / SeedAdmin
   ├─ corsMiddleware        CORS + preflight OPTIONS
   ├─ http.ServeMux         роутинг "GET /user" (Go 1.22)
   └─ HomeHandler / HealthHandler

internal/
   ├─ httpx/                ParseID, WriteJSON, WriteError
   ├─ errs/                 ErrNotFound, ErrUnauthorized
   ├─ auth/                 HashPassword, CheckPassword, HashOptional, ValidatePassword
   ├─ session/              Create, Validate, Revoke (SHA-256 токенов в БД)
   ├─ database/
   │    ├─ database.go      подключение с retry
   │    ├─ migrations.go    golang-migrate + migrations/*.sql
   │    └─ seed.go          роли admin/user и технический пользователь
   └─ <entity>/
        ├─ models.go        структуры запроса/ответа
        ├─ repository.go    SQL (Query/QueryRow/Exec)
        └─ controller/      HTTP-обработчики + Handler
```

Слои: `controller` → `repository` → PostgreSQL. Контроллер отвечает за HTTP и
коды ответов, репозиторий — за SQL и маппинг строк в модели. Аутентификация
живёт в `user/controller` и использует пакеты `auth` и `session`.
`session.Validate` подготовлен для будущего `authMiddleware`, но пока нигде не
вызывается.

## Структура проекта

```
backend-go/
├─ main.go                 точка входа, маршруты, CORS, раздача статики
├─ web/                    статические страницы (go:embed)
│  ├─ index.html           витрина товаров
│  ├─ auth.html            вход и регистрация
│  └─ admin_panel.html     CRUD-панель (+ /health)
├─ internal/
│  ├─ auth/                bcrypt и валидация пароля
│  ├─ database/            InitDB, CloseDB, RunMigrations, SeedAdmin, migrations/
│  ├─ errs/                общие ошибки
│  ├─ httpx/               HTTP/JSON-хелперы
│  ├─ product/             товары
│  ├─ rules/               роли
│  ├─ session/             серверные сессии
│  ├─ status/              статусы товаров
│  ├─ user/                пользователи, login/register/logout
│  └─ user_status/         статусы пользователей
├─ docs/                   документация
├─ Dockerfile              multi-stage сборка
└─ docker-compose.yml      app + postgres
```

Каждая сущность устроена одинаково: `models.go`, `repository.go`,
`controller/{handler,create,get,get_all,update,delete}.go`. У `user` есть
дополнительные файлы `login.go`, `register.go`, `logout.go`.

## Сущности

| Сущность            | Таблица       | Связи / назначение                            |
|---------------------|---------------|-----------------------------------------------|
| Пользователь        | `users`       | `status_id → user_status`, `rule_id → rules`  |
| Статус пользователя | `user_status` | справочник                                    |
| Роль                | `rules`       | справочник (`admin`, `user`)                  |
| Статус товара       | `status`      | справочник                                    |
| Товар               | `products`    | `status_id → status`                          |
| Сессия              | `sessions`    | `user_id → users` (не CRUD-сущность)          |

Подробная схема — `docs/database/schema.md`, картинка — `docs/ERD.png`.
Актуальный перечень колонок и типов — в SQL-миграциях
`internal/database/migrations/`.

## API

`id` передаётся query-параметром, например `GET /user?id=1`.

| Метод  | Путь                      | Описание                                    |
|--------|---------------------------|---------------------------------------------|
| GET    | `/`                       | витрина (`web/index.html`)                  |
| GET    | `/auth.html`              | вход/регистрация                            |
| GET    | `/admin_panel.html`       | админ-панель                                |
| GET    | `/hello`                  | приветствие (старый `/`)                    |
| GET    | `/health`                 | ping БД                                     |
| POST   | `/user/register`          | регистрация, роль `user`                    |
| POST   | `/user/login`             | вход, ответ `{"ok": true}`, кука `session`  |
| POST   | `/user/logout`            | выход по куке `session`                     |
| GET    | `/user/me`                | текущий пользователь по куке `session`      |
| GET    | `/users`                  | список пользователей                        |
| GET    | `/user?id=`               | пользователь по id                          |
| POST   | `/user/create`            | создать пользователя                        |
| PUT    | `/user/update?id=`        | обновить пользователя                       |
| DELETE | `/user/delete?id=`        | удалить пользователя                        |
| GET    | `/statuses`               | список статусов товаров                     |
| GET    | `/status?id=`             | статус товара                               |
| POST   | `/status/create`          | создать статус                              |
| PUT    | `/status/update?id=`      | обновить статус                             |
| DELETE | `/status/delete?id=`      | удалить статус                              |
| GET    | `/rules`                  | список ролей                                |
| GET    | `/rule?id=`               | роль по id                                  |
| POST   | `/rule/create`            | создать роль                                |
| PUT    | `/rule/update?id=`        | обновить роль                               |
| DELETE | `/rule/delete?id=`        | удалить роль                                |
| GET    | `/user-statuses`          | список статусов пользователей               |
| GET    | `/user-status?id=`        | статус пользователя                         |
| POST   | `/user-status/create`     | создать статус                              |
| PUT    | `/user-status/update?id=` | обновить статус                             |
| DELETE | `/user-status/delete?id=` | удалить статус                              |
| GET    | `/products`               | список товаров                              |
| GET    | `/product?id=`            | товар по id                                 |
| POST   | `/product/create`         | создать товар                               |
| PUT    | `/product/update?id=`     | обновить товар                              |
| DELETE | `/product/delete?id=`     | удалить товар                               |

Коды ответов: `200` — `GET`/`PUT` и login, `201` — create и register, `204` —
delete и logout, `400` — неверный JSON/`id`/обязательные поля/слабый пароль,
`401` — неверные учётные данные или отсутствующий токен, `404` — сущность не
найдена, `409` — email уже зарегистрирован, `500` — внутренняя ошибка.
Ошибки отдаются в формате `{"error": "..."}`.

Админ-панель проверяет длину пароля (8–72 байта, как и бэкенд) и лимиты
строковых полей из схемы БД до отправки запроса. При `401`/`403` она очищает
локальную сессию и возвращает пользователя на `auth.html`.

Формат ответа пользователя (статус и роль развёрнуты, `null` если не назначены;
`created_at` / `updated_at` в RFC3339 и видны в таблице панели; `password` в
ответ не попадает):

```json
{
  "id": 1,
  "name": "Vadim",
  "surname": "Garaev",
  "female": null,
  "email": "user@example.com",
  "created_at": "2026-11-19T12:34:56.789Z",
  "updated_at": "2026-11-19T12:34:56.789Z",
  "status": { "id": 1, "name": "Активен", "code": "active" },
  "rule": { "id": 1, "name": "Admin", "code": "admin" }
}
```

## Переменные окружения

| Переменная       | Значение по умолчанию | Назначение                     |
|------------------|-----------------------|--------------------------------|
| `DB_HOST`        | —                     | хост PostgreSQL                |
| `DB_PORT`        | —                     | порт PostgreSQL                |
| `DB_USER`        | —                     | пользователь БД                |
| `DB_PASSWORD`    | —                     | пароль БД                      |
| `DB_NAME`        | —                     | имя БД                         |
| `ADMIN_EMAIL`    | `email@email`         | email технического админа      |
| `ADMIN_PASSWORD` | `1235`                | пароль технического админа     |
| `ADMIN_NAME`     | `tech`                | имя технического админа        |

При старте приложение применяет миграции и через `SeedAdmin` создаёт роли
`admin`/`user` и технического пользователя с bcrypt-хешем пароля. Если пароль
в БД отличается от `ADMIN_PASSWORD`, он перезаписывается.

## Запуск

```bash
docker compose up --build
```

Приложение и страницы: `http://localhost:8080`, PostgreSQL: `localhost:5432`, БД `myapp`. Открывать HTML через `file://` не нужно — cookie-сессия и same-origin API-запросы работают только через HTTP.

Локально без Docker:

```bash
export DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=myapp
go run .
```

## Документация
 
| Файл                          | О чём                         |
|-------------------------------|-------------------------------|
| `docs/project/implemented.md` | что уже реализовано           |
| `docs/project/tasks.md`       | задачи на будущее             |
| `docs/database/schema.md`     | схема БД и связи              |
| `docs/go/`                    | шпаргалки по Go               |
| `docs/docker/`                | шпаргалки по Docker           |
