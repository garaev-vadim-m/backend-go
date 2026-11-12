# Что уже реализовано

Отметки: `[x]` — готово, `[ ]` — не сделано (см. `tasks.md`).

## Инфраструктура

- [x] Go-модуль `hello-go-backend`, Go 1.22
- [x] Подключение к PostgreSQL через `database/sql` + `lib/pq`
- [x] Retry подключения: 5 попыток с паузой 2 с (`internal/database/database.go`)
- [x] Миграции при старте: `CREATE TABLE IF NOT EXISTS` для 5 таблиц (`internal/database/migrations.go`)
- [x] Dockerfile с multi-stage сборкой (builder → alpine)
- [x] `docker-compose.yml`: app + `postgres:15-alpine`, healthcheck, named volume
- [x] `.dockerignore`
- [x] Endpoint `/health` с ping БД
- [x] CORS-middleware и preflight `OPTIONS` (`main.go`)

## Модель данных

- [x] Таблицы: `users`, `user_status`, `rules`, `status`, `products`
- [x] Внешние ключи: `users.status_id`, `users.rule_id`, `products.status_id`
- [x] Nullable-поля: `users.female`, `users.status_id`, `users.rule_id`,
      `products.description`, `products.status_id`
- [x] ERD: `docs/ERD.png`, описание — `docs/database/schema.md`

## API и CRUD

- [x] Полный CRUD для `users` — `/users`, `/user`, `/user/{create,update,delete}`
- [x] Полный CRUD для `status` — `/statuses`, `/status`, `/status/{create,update,delete}`
- [x] Полный CRUD для `rules` — `/rules`, `/rule`, `/rule/{create,update,delete}`
- [x] Полный CRUD для `user_status` — `/user-statuses`, `/user-status`, `/user-status/{create,update,delete}`
- [x] Полный CRUD для `products` — `/products`, `/product`, `/product/{create,update,delete}`
- [x] Роутинг с методами Go 1.22 (`"GET /user"`)
- [x] `POST .../create` возвращает `201 Created`
- [x] `GET` / `PUT` возвращают `200 OK`
- [x] `DELETE` возвращает `204 No Content`

## Паттерны и качество кода

- [x] Единый пакет ошибок `internal/errs` с `ErrNotFound`
- [x] Обёртка ошибок через `%w` и проверка `errors.Is`
- [x] HTTP/JSON-хелперы `internal/httpx`: `ParseID`, `WriteJSON`, `WriteError`
- [x] Единый формат ошибки `{"error": "..."}`
- [x] Handler-структура с `*sql.DB` и конструктором `NewHandler`
- [x] Общий `scanner`-интерфейс для `*sql.Row` и `*sql.Rows`
- [x] JOIN статусов/ролей в выборках (развёрнутые объекты вместо `*_id`)
- [x] `INSERT` / `UPDATE` через CTE + `RETURNING` + `JOIN`
- [x] `sql.Null*` для nullable-полей из JOIN
- [x] Проверка `id` в query (`httpx.ParseID`)
- [x] Проверка обязательного `name` при создании/обновлении
- [x] Проверка `RowsAffected` при удалении
- [x] Логирование внутренних ошибок через `log.Printf`

## Веб-панель

- [x] `index.html`: CRUD для `users`, `status`, `rules`, `user_status`, `products`
- [x] Индикатор `/health`
- [x] Настройка базового URL API (сохраняется в `localStorage`)

## Документация

- [x] `README.md`
- [x] `docs/go/` — connection, fmt, pointers, syntax, errors, interfaces, sql, json, http
- [x] `docs/docker/` — basics, build-run, compose
- [x] `docs/database/schema.md`
- [x] `docs/План действий.md`
- [x] `docs/project/` — overview, implemented, tasks

## Чего пока нет

- [ ] Аутентификации и авторизации
- [ ] Ограничения прав по ролям (admin)
- [ ] Корзины и Redis
- [ ] Тестов (unit/integration)
- [ ] Пагинации, фильтрации, сортировки
- [ ] OpenAPI/Swagger
- [ ] Валидации внешних ключей (`status_id`, `rule_id`) с кодом 400/404
- [ ] Graceful shutdown и структурированного логирования
