# Схема базы данных

## ERD

Визуальная схема базовых таблиц лежит в `docs/ERD.png`. Актуальный источник
истины по колонкам и типам — SQL-файлы в `internal/database/migrations/`.
Текстовый вариант:

```
 ┌─────────────────┐              ┌─────────────────┐
 │   user_status   │              │      rules      │
 ├─────────────────┤              ├─────────────────┤
 │ id  PK          │              │ id  PK          │
 │ name            │              │ name            │
 │ code            │              │ code            │
 └────────┬────────┘              └────────┬────────┘
          │ status_id                      │ rule_id
          └───────────────┬────────────────┘
                          ▼
                 ┌─────────────────┐
                 │      users      │
                 ├─────────────────┤
                 │ id  PK          │
                 │ name            │
                 │ surname  (null) │
                 │ female   (null) │
                 │ email    (UQ)   │
                 │ password (bcrypt)│
                 │ status_id  FK   │
                 │ rule_id    FK   │
                 └────────┬────────┘
                          │ user_id (ON DELETE CASCADE)
                          ▼
                 ┌─────────────────┐
                 │     sessions    │
                 ├─────────────────┤
                 │ id  PK          │
                 │ user_id   FK    │
                 │ token_hash UQ   │
                 │ created_at      │
                 │ expires_at      │
                 │ revoked_at      │
                 └─────────────────┘

 ┌─────────────────┐              ┌─────────────────┐
 │      status     │              │     products    │
 ├─────────────────┤              ├─────────────────┤
 │ id  PK          │◄─────────────┤ id  PK          │
 │ name            │ status_id FK │ name            │
 │ code            │              │ description     │
 └─────────────────┘              │ status_id  FK   │
                                  └─────────────────┘
```

## Связи

| Откуда         | Поле        | Куда          | Поведение при удалении |
|----------------|-------------|---------------|------------------------|
| `users`        | `status_id` | `user_status` | по умолчанию (`NO ACTION`) |
| `users`        | `rule_id`   | `rules`       | по умолчанию (`NO ACTION`) |
| `products`     | `status_id` | `status`      | по умолчанию (`NO ACTION`) |
| `sessions`     | `user_id`   | `users`       | `ON DELETE CASCADE`        |

Обратите внимание: у пользователя статус берётся из `user_status`, а у продукта — из `status`.

## Таблицы

| Таблица       | Колонки                                                          |
|---------------|------------------------------------------------------------------|
| `status`      | `id` (PK), `name`, `code`                                        |
| `rules`       | `id` (PK), `name`, `code`                                        |
| `user_status` | `id` (PK), `name`, `code`                                        |
| `users`       | `id` (PK), `name`, `surname` (nullable), `female` (nullable), `email` (nullable, UNIQUE), `password` (nullable, bcrypt), `status_id` (FK), `rule_id` (FK) |
| `products`    | `id` (PK), `name`, `description` (nullable), `status_id` (FK)    |
| `sessions`    | `id` (PK), `user_id` (FK, NOT NULL), `token_hash` (UNIQUE), `created_at`, `expires_at`, `revoked_at` (nullable) |

`id` — `SERIAL PRIMARY KEY`, то есть автоинкремент. `name` и `code` — `VARCHAR(255) NOT NULL`.
`female` — `VARCHAR(10)`, остальные строковые поля — `VARCHAR(255)`.
`created_at` и `expires_at` в `sessions` — `TIMESTAMPTZ NOT NULL`
(`created_at` по умолчанию `now()`), `revoked_at` — `TIMESTAMPTZ NULL`.

`users.email` формально nullable (PostgreSQL разрешает несколько `NULL`), но
регистрация и сид админа всегда его заполняют; уникальность даёт ограничение
`users_email_key` из миграции 000009. `users.password` тоже nullable — у
пользователей без пароля вход невозможен.

В `sessions` хранится не сам токен, а его SHA-256-хеш (64 hex-символа).
Токен выдаётся клиенту при `/user/login` и живёт 24 часа; `revoked_at`
используется для выхода. Индекс `sessions_user_id_idx` ускоряет поиск сессий
пользователя.

## Миграции

Миграции применяются при старте `database.RunMigrations()` через
`github.com/golang-migrate/migrate/v4`. SQL-файлы встроены в бинарник через
`//go:embed migrations/*.sql`; отдельное соединение собирается из `DB_*` в
`databaseURL()`.

| Версия | Файл                              | Что делает                              |
|--------|-----------------------------------|-----------------------------------------|
| 000001 | `000001_create_status.up.sql`     | таблица `status`                        |
| 000002 | `000002_create_rules.up.sql`      | таблица `rules`                         |
| 000003 | `000003_create_user_status.up.sql`| таблица `user_status`                   |
| 000004 | `000004_create_users.up.sql`      | таблица `users`                         |
| 000005 | `000005_create_products.up.sql`   | таблица `products`                      |
| 000006 | `000006_add_password_to_users.up.sql` | `users.password`                   |
| 000007 | `000007_create_sessions.up.sql`   | таблица `sessions`                      |
| 000008 | `000008_add_revoked_at_sessions.up.sql` | `revoked_at` + индекс             |
| 000009 | `000009_unique_users_email.up.sql`| уникальный `users.email`                |

Порядок важен: сначала таблицы без зависимостей, потом `users`, `products` и
`sessions`, которые ссылаются на них через `REFERENCES`. Для миграций
000006–000009 есть `.down.sql`; у 000001–000005 down-файлов пока нет.

Роли `admin`/`user` и технический пользователь создаются не миграцией, а
`database.SeedAdmin()` при каждом старте приложения.

## Схема и пакеты

| Таблица / тема | Go-пакет / файл            |
|----------------|----------------------------|
| `status`       | `internal/status`          |
| `rules`        | `internal/rules`           |
| `user_status`  | `internal/user_status`     |
| `users`        | `internal/user`            |
| `products`     | `internal/product`         |
| `sessions`     | `internal/session`         |
| Пароли         | `internal/auth`            |
| Миграции, сид  | `internal/database`        |

В пакетах сущностей: `models.go` (структуры), `repository.go` (SQL),
`controller/` (HTTP-обработчики).

---

## Чек-лист

- [x] таблицы-справочники создаются раньше зависимых
- [x] связи описаны через `REFERENCES`
- [x] nullable-колонки в Go представлены указателями или `sql.Null*`
- [x] текстовое описание соответствует SQL-миграциям
