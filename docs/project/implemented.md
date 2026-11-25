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
- [x] Статика (`web/*.html`) встроена в бинарник через `embed.FS` и раздаётся
      тем же сервером: `/`, `/auth.html`, `/admin_panel.html`; приветствие
      перенесено на `/hello` (`main.go`)
- [x] Endpoint `/health` с ping БД
- [x] CORS-middleware и preflight `OPTIONS` (`main.go`), заголовок
      `Authorization` разрешён; `Access-Control-Allow-Origin` отражает `Origin`

## Модель данных

- [x] Таблицы: `users`, `user_status`, `rules`, `status`, `products`, `sessions`
- [x] Миграции: 12 версий в `internal/database/migrations/` (создание таблиц,
      `users.password`, `sessions`, `sessions.revoked_at` + индекс,
      уникальность `users.email`, `created_at`/`updated_at`, индекс `expires_at`)
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
- [x] TTL сессии 24 часа; `Validate` учитывает `expires_at` и `revoked_at` и
      вызывается в `/user/me` и в `middleware.Auth.Require`
- [x] Продление сессии: `Renew` сдвигает `expires_at`, если до конца меньше 12 ч;
      кука перевыпускается при куки-авторизации
- [x] Logout отзывает сессию через `revoked_at`, строка остаётся в истории
- [x] `errs.ErrUnauthorized` для невалидной сессии
- [x] `GET /user/me`, `Auth.Require`, `Auth.RequireRole("admin")` —
      защищённые роуты, id пользователя в `context.Context`
- [x] Кука `session`: `HttpOnly`, `SameSite=Lax`, `Secure` при `COOKIE_SECURE=true`

## Паттерны и качество кода

- [x] Единый пакет ошибок `internal/errs` с `ErrNotFound` и `ErrUnauthorized`
- [x] Обёртка ошибок через `%w` и проверка `errors.Is` (не везде: 4 репозитория
      UPDATE не оборачивают `ErrNotFound`)
- [x] HTTP/JSON-хелперы `internal/httpx`: `ParseID`, `WriteJSON`, `WriteError`
- [x] Формат ошибки `{"error": "..."}` (кроме текстовых `/health` и списков
      при ошибке БД)
- [x] Handler-структура с `*sql.DB` и конструктором `NewHandler`
- [x] Общий `scanner`-интерфейс для `*sql.Row` и `*sql.Rows`
- [x] JOIN статусов/ролей в выборках (развёрнутые объекты вместо `*_id`)
- [x] `users.created_at` / `users.updated_at` (миграции 000010/000011) попадают
      в ответ; `updated_at` пересчитывается в `UPDATE`
- [x] `INSERT` / `UPDATE` через CTE + `RETURNING` + `JOIN`
- [x] `sql.Null*` для nullable-полей из JOIN
- [x] Проверка `id` в query (`httpx.ParseID`)
- [x] Проверка обязательного `name` при создании/обновлении
- [x] Проверка пароля при создании/обновлении пользователя; пустой пароль
      при `PUT` оставляет старый (`COALESCE(NULLIF($5, ''), password)`)
- [x] Проверка `RowsAffected` при удалении
- [x] Логирование внутренних ошибок через `log.Printf`

## Веб-страницы

- [x] `web/index.html`: витрина товаров из `/products`, ссылка на админ-панель
      для роли `admin`, индикатор сессии
- [x] `web/auth.html`: вход и регистрация; серверная сессия живёт в
      `HttpOnly`-куке `session`, в `localStorage` (`hello-go-backend.auth`)
      сохраняется только профиль для UI, админ перенаправляется в
      `/admin_panel.html`
- [x] `web/admin_panel.html`: CRUD для `users`, `status`, `rules`, `user_status`,
      `products`, проверка `/health`, выход из сессии
- [x] Текущий пользователь и роль во всех страницах берутся из `GET /user/me`
- [x] Базовый URL API по умолчанию пустой (текущий origin), явное значение
      сохраняется в `localStorage`; запросы идут с `credentials: "include"`,
      `Authorization` не используется
- [x] Кнопки создания/изменения/удаления в панели доступны только роли `admin`;
      сервер дублирует проверку через `Auth.RequireRole("admin")`
- [x] Пароль в форме пользователей проверяется по правилам бэкенда: 8–72 байта
      (`minlength`/`maxlength` + проверка `TextEncoder` на отправке); пустой
      пароль при редактировании оставляет прежний
- [x] Строковые поля ограничены длинами из схемы БД: 255 символов, `female` — 10
- [x] Таблица пользователей показывает `created_at` / `updated_at` в локальном
      формате; карточкой по ID эти поля тоже видны в JSON
- [x] При `401` обработчик `api()` очищает локальную сессию и возвращает на
      `auth.html`; при `403` показывает ошибку прав без выхода; logout также
      идёт через этот обработчик
- [x] Заголовок панели показывает фактическую роль из сессии, а не всегда `admin`;
      атрибут `hidden` гарантированно скрывает элементы

## Документация

- [x] `README.md`
- [x] `docs/go/` — connection, fmt, pointers, syntax, errors, interfaces, sql, json, http
- [x] `docs/docker/` — basics, build-run, compose
- [x] `docs/database/schema.md`, `docs/ERD.png`
- [x] `docs/project/` — overview, implemented, tasks

## Чего пока нет

- [ ] Перевода всех пользователей с открытых паролей на bcrypt (fallback в
      `login.go` ещё есть)
- [ ] JWT и refresh-токенов
- [ ] Корзины и Redis
- [ ] Тестов (unit/integration)
- [ ] Пагинации, фильтрации, сортировки
- [ ] OpenAPI/Swagger
- [ ] Валидации внешних ключей (`status_id`, `rule_id`) с кодом 400/404
- [ ] Graceful shutdown и структурированного логирования

## Известные ограничения

- [ ] `PUT` для несуществующего `product`/`status`/`rule`/`user_status`
      возвращает `500` вместо `404`: репозитории не оборачивают `ErrNotFound`
- [ ] Дублирующийся email в `/user/create` даёт `500`; `409` обрабатывается только
      в `/user/register`
- [ ] Списки при ошибке БД отдают текст `Error: ...`, а не `{"error": "..."}`
- [ ] `/user/me` не принимает `Authorization: Bearer` и не продлевает сессию
- [ ] `Secure` куки зависит от `COOKIE_SECURE` (по умолчанию `false`)
