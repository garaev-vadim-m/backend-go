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

## Формат пользователя

`GET /users`, `GET /user?id=`, `POST /user/create` и `PUT /user/update?id=` отдают
одинаковый JSON — статус и роль уже развёрнуты в объекты:

```json
{
  "id": 1,
  "name": "Vadim",
  "female": "Garaev",
  "status": { "id": 1, "name": "Активен", "code": "active" },
  "rule": { "id": 1, "name": "Admin", "code": "admin" }
}
```

`status` и `rule` равны `null`, если не назначены. Идентификаторы `status_id` и
`rule_id` передаются только в теле запросов `POST`/`PUT` (`null` = не назначено).

## Веб-панель (`index.html`)

В корне есть `index.html` — самодостаточная страница (чистый JS, без сборки) для
работы с API: CRUD для `users`, `status`, `rules`, `user_status`, плюс индикатор
`/health`. Откройте файл в браузере (двойным щелчком или через любой статический
сервер) и укажите базовый URL API в поле сверху (по умолчанию
`http://localhost:8080`; значение сохраняется в `localStorage`).

Запросы из браузера идут с другого origin (`file://` или другой порт), поэтому
бэкенд обёрнут в CORS-middleware (`corsMiddleware` в `main.go`): он отдаёт
`Access-Control-Allow-Origin/-Methods/-Headers` и отвечает `204` на preflight
`OPTIONS`. Если в панели появляется «Сетевая/CORS ошибка» — проверьте, что
бэкенд запущен и доступен по указанному URL.

## Локальная разработка (без Docker)

1. Запустить PostgreSQL и создать БД `myapp`
2. Экспортировать переменные окружения:

```bash
export DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=myapp
```

Подключение к базе данных
```bash
docker-compose exec postgres psql -U postgres -d myapp
```

Проверка таблиц БД
```bash
docker-compose exec postgres psql -U postgres -d myapp -c "\dt"
```

3. Запустить:

```bash
go run .
```

## Стек

- Go 1.27
- PostgreSQL 15
- Docker / Docker Compose
