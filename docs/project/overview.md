# Проект hello-go-backend

## Что это

REST API на Go с PostgreSQL и одностраничной веб-панелью. Проект реализует CRUD
для пяти сущностей, связи между ними и единый формат JSON-ответов.

- API: `http://localhost:8080`
- Веб-панель: `index.html` (чистый JS, без сборки)
- Документация: `docs/`

## Стек

| Компонент  | Версия / технология               |
|------------|-----------------------------------|
| Язык       | Go 1.22                           |
| БД         | PostgreSQL 15                     |
| Драйвер БД | `github.com/lib/pq`               |
| Контейнеры | Docker, Docker Compose            |
| Фронтенд   | HTML + vanilla JS (`index.html`)  |

## Архитектура

```
main.go
  ├─ corsMiddleware        CORS + preflight OPTIONS
  ├─ http.ServeMux         роутинг "GET /user" (Go 1.22)
  └─ database.InitDB / RunMigrations

internal/
  ├─ httpx/                ParseID, WriteJSON, WriteError
  ├─ errs/                 ErrNotFound (sentinel-ошибка)
  ├─ database/             подключение + миграции
  └─ <entity>/
       ├─ models.go        структуры запроса/ответа
       ├─ repository.go    SQL (Query/QueryRow/Exec)
       └─ controller/      HTTP-обработчики + Handler
```

Слои: `controller` → `repository` → PostgreSQL. Контроллер отвечает за HTTP и
коды ответов, репозиторий — за SQL и маппинг строк в модели.

## Структура проекта

```
backend-go/
├─ main.go                 точка входа, маршруты, CORS
├─ index.html              веб-панель (CRUD + /health)
├─ internal/
│  ├─ database/            InitDB, CloseDB, RunMigrations
│  ├─ errs/                общие ошибки
│  ├─ httpx/               HTTP/JSON-хелперы
│  ├─ product/             товары
│  ├─ rules/               роли
│  ├─ status/              статусы товаров
│  ├─ user/                пользователи
│  └─ user_status/         статусы пользователей
├─ docs/                   документация
├─ Dockerfile              multi-stage сборка
└─ docker-compose.yml      app + postgres
```

Каждая сущность устроена одинаково: `models.go`, `repository.go`,
`controller/{handler,create,get,get_all,update,delete}.go`.

## Сущности

| Сущность          | Таблица       | Связи                                        |
|-------------------|---------------|----------------------------------------------|
| Пользователь      | `users`       | `status_id → user_status`, `rule_id → rules` |
| Статус пользователя | `user_status` | справочник                                 |
| Роль              | `rules`       | справочник                                   |
| Статус товара     | `status`      | справочник                                   |
| Товар             | `products`    | `status_id → status`                         |

Подробная схема — `docs/database/schema.md`, картинка — `docs/ERD.png`.

## API

| Метод  | Путь                      | Описание                        |
|--------|---------------------------|---------------------------------|
| GET    | `/`                       | приветствие                     |
| GET    | `/health`                 | ping БД                         |
| GET    | `/users`                  | список пользователей            |
| GET    | `/user?id=`               | пользователь по id              |
| POST   | `/user/create`            | создать пользователя            |
| PUT    | `/user/update?id=`        | обновить пользователя           |
| DELETE | `/user/delete?id=`        | удалить пользователя            |
| GET    | `/statuses`               | список статусов товаров         |
| GET    | `/status?id=`             | статус товара                   |
| POST   | `/status/create`          | создать статус                  |
| PUT    | `/status/update?id=`      | обновить статус                 |
| DELETE | `/status/delete?id=`      | удалить статус                  |
| GET    | `/rules`                  | список ролей                    |
| GET    | `/rule?id=`               | роль по id                      |
| POST   | `/rule/create`            | создать роль                    |
| PUT    | `/rule/update?id=`        | обновить роль                   |
| DELETE | `/rule/delete?id=`        | удалить роль                    |
| GET    | `/user-statuses`          | список статусов пользователей   |
| GET    | `/user-status?id=`        | статус пользователя             |
| POST   | `/user-status/create`     | создать статус                  |
| PUT    | `/user-status/update?id=` | обновить статус                 |
| DELETE | `/user-status/delete?id=` | удалить статус                  |
| GET    | `/products`               | список товаров                  |
| GET    | `/product?id=`            | товар по id                     |
| POST   | `/product/create`         | создать товар                   |
| PUT    | `/product/update?id=`     | обновить товар                  |
| DELETE | `/product/delete?id=`     | удалить товар                   |

Формат ответа пользователя (статус и роль развёрнуты, `null` если не назначены):

```json
{
  "id": 1,
  "name": "Vadim",
  "female": "Garaev",
  "status": { "id": 1, "name": "Активен", "code": "active" },
  "rule": { "id": 1, "name": "Admin", "code": "admin" }
}
```

## Запуск

```bash
docker compose up --build
```

Приложение: `http://localhost:8080`, PostgreSQL: `localhost:5432`, БД `myapp`.

Локально без Docker:

```bash
export DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=myapp
go run .
```

## Документация

| Файл                          | О чём                         |
|-------------------------------|-------------------------------|
| `docs/План действий.md`       | план развития проекта         |
| `docs/project/implemented.md` | что уже реализовано           |
| `docs/project/tasks.md`       | задачи на будущее             |
| `docs/database/schema.md`     | схема БД и связи              |
| `docs/go/`                    | шпаргалки по Go               |
| `docs/docker/`                | шпаргалки по Docker           |
