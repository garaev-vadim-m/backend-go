# Что уже реализовано

Отметки: `[x]` — готово, `[ ]` — не сделано (см. `tasks.md`).

## Инфраструктура

- [x] Go-модуль `hello-go-backend`, Go 1.22
- [x] Подключение к PostgreSQL через `database/sql` + `lib/pq`
- [x] Retry подключения: 5 попыток с паузой 2 с (`internal/database/database.go`)
- [x] Миграции через `golang-migrate` v4: SQL-файлы встроены в бинарник через
      `embed.FS` и применяются при старте (`internal/database/migrations.go`)
- [x] Сидирование при старте: роли `admin`/`user` и технический пользователь
      с bcrypt-паролем (`internal/database/seed.go`, env `ADMIN_*`)
- [x] Dockerfile с multi-stage сборкой (builder → alpine)
- [x] `docker-compose.yml`: app + `postgres:15-alpine`, healthcheck, named volume
- [x] `.dockerignore`
- [x] Endpoint `/health` с ping БД
- [x] CORS-middleware и preflight `OPTIONS` (`main.go`), заголовок
      `Authorization` разрешён; `Access-Control-Allow-Origin` отражает `Origin`

## Модель данных

- [x] Таблицы: `users`, `user_status`, `rules`, `status`, `products`, `sessions`
- [x] Миграции: 9 версий в `internal/database/migrations/` (создание таблиц,
      `users.password`, `sessions`, `sessions.revoked_at` + индекс,
      уникальность `users.email`)
- [x] Внешние ключи: `users.status_id`, `users.rule_id`, `products.status_id`,
      `sessions.user_id` (`ON DELETE CASCADE`)
- [x] Nullable-поля: `users.surname`, `users.female`, `users.email`,
      `users.password`, `users.status_id`, `users.rule_id`,
      `products.description`, `products.status_id`, `sessions.revoked_at`
- [x] `users.email` уникален (PostgreSQL допускает несколько `NULL`)
- [x] ERD: `docs/ERD.png`, текстовое описание — `docs/database/schema.md`

## API и CRUD

- [x] Полный CRUD для `users` — `/users`, `/user`, `/user/{create,update,delete}`
- [x] Полный CRUD для `status` — `/statuses`, `/status`, `/status/{create,update,delete}`
- [x] Полный CRUD для `rules` — `/rules`, `/rule`, `/rule/{create,update,delete}`
- [x] Полный CRUD для `user_status` — `/user-statuses`, `/user-status`, `/user-status/{create,update,delete}`
- [x] Полный CRUD для `products` — `/products`, `/product`, `/product/{create,update,delete}`
- [x] Регистрация `POST /user/register`, вход `POST /user/login`,
      выход `POST /user/logout`
- [x] Роутинг с методами Go 1.22 (`"GET /user"`)
- [x] `POST .../create` и register возвращают `201 Created`
- [x] `GET` / `PUT` возвращают `200 OK`
- [x] `DELETE` и logout возвращают `204 No Content`
- [x] `401` при неверных учётных данных и при logout без токена

## Аутентификация

- [x] Пароли хранятся как bcrypt-хеши (`internal/auth`), длина пароля 8–72 байт
- [x] При входе поддерживается старый пароль в открытом виде (сравнение
      `subtle.ConstantTimeCompare`), после входа выдаётся новая сессия
- [x] Регистрация не принимает `rule_id`/`status_id`: роль `user` назначается
      через `rules.GetIDByCode`
- [x] Дублирующийся email даёт `409 Conflict` (`pq.Error` код `23505`)
- [x] Сессия: случайный токен 32 байта (hex), в БД хранится только его
      SHA-256-хеш (`internal/session`)
- [x] TTL сессии 24 часа; `Validate` учитывает `expires_at` и `revoked_at`
      (пока не вызывается роутами)
- [x] Logout отзывает сессию через `revoked_at`, строка остаётся в истории
- [x] `errs.ErrUnauthorized` для невалидной сессии
- [ ] `/auth/me`, `authMiddleware`, `requireRole("admin")` — см. `tasks.md`

## Паттерны и качество кода

- [x] Единый пакет ошибок `internal/errs` с `ErrNotFound` и `ErrUnauthorized`
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
- [x] Проверка пароля при создании/обновлении пользователя; пустой пароль
      при `PUT` оставляет старый (`COALESCE(NULLIF($5, ''), password)`)
- [x] Проверка `RowsAffected` при удалении
- [x] Логирование внутренних ошибок через `log.Printf`

## Веб-страницы

- [x] `index.html`: витрина товаров из `/products`, ссылка на админ-панель для
      роли `admin`, индикатор сессии
- [x] `auth.html`: вход и регистрация, сохранение `{token, user}` в
      `localStorage` (`hello-go-backend.auth`), перенаправляет админа в `admin_panel.html`
- [x] `admin_panel.html`: CRUD для `users`, `status`, `rules`, `user_status`,
      `products`, проверка `/health`, выход из сессии
- [x] Настройка базового URL API (сохраняется в `localStorage`), заголовок
      `Authorization: Bearer <token>` добавляется в запросы
- [x] Кнопки создания/изменения/удаления в панели доступны только роли `admin`
      (проверка только на клиенте; сервер её не применяет)

## Документация

- [x] `README.md`
- [x] `docs/go/` — connection, fmt, pointers, syntax, errors, interfaces, sql, json, http
- [x] `docs/docker/` — basics, build-run, compose
- [x] `docs/database/schema.md`, `docs/ERD.png`
- [x] `docs/project/` — overview, implemented, tasks

## Чего пока нет

- [ ] Ограничения прав по ролям на стороне сервера (проверка токена и `admin`
      на CRUD-роутах)
- [ ] Защищённых маршрутов `/auth/me`, middleware и хранения пользователя в
      `context.Context`
- [ ] TTL-продления сессии, cookie с `HttpOnly`/`SameSite`/`Secure`
- [ ] JWT и refresh-токенов
- [ ] Колонок `created_at` / `updated_at` у пользователей
- [ ] Корзины и Redis
- [ ] Тестов (unit/integration)
- [ ] Пагинации, фильтрации, сортировки
- [ ] OpenAPI/Swagger
- [ ] Валидации внешних ключей (`status_id`, `rule_id`) с кодом 400/404
- [ ] Graceful shutdown и структурированного логирования
