# Схема базы данных

## ERD

Визуальная схема базовых таблиц лежит в `docs/ERD.png` (картинка не включает
корзину и заказы). Актуальный источник истины по колонкам и типам — SQL-файлы
в `internal/database/migrations/`. Текстовый вариант:

```
user_status (id, name, code) ─┐
                              ├─ users (id, name, email, password,
rules       (id, name, code) ─┘         status_id FK, rule_id FK,
                                        created_at, updated_at)
                                          ├─ sessions (id, user_id FK,
                                          │   token_hash UQ, created_at,
                                          │   expires_at, revoked_at)
                                          ├─ carts (id, user_id FK/UQ,
                                          │   created_at, updated_at)
                                          │    └─ cart_items (cart_id FK,
                                          │         product_id FK, quantity)
                                          └─ orders (id, user_id FK, status,
                                              total, idempotency_key UQ
                                              по паре с user_id, created_at)
                                               └─ order_items (order_id FK,
                                                    product_id FK, quantity,
                                                    price — снимок цены)

status (id, name, code) ─ products (id, name, description, status_id FK,
                                    price >= 0, stock >= 0)
                              ├─ cart_items.product_id FK
                              └─ order_items.product_id FK
```

## Связи

| Откуда         | Поле        | Куда          | Поведение при удалении |
|----------------|-------------|---------------|------------------------|
| `users`        | `status_id` | `user_status` | по умолчанию (`NO ACTION`) |
| `users`        | `rule_id`   | `rules`       | по умолчанию (`NO ACTION`) |
| `products`     | `status_id` | `status`      | по умолчанию (`NO ACTION`) |
| `sessions`     | `user_id`   | `users`       | `ON DELETE CASCADE`        |
| `carts`        | `user_id`   | `users`       | `ON DELETE CASCADE`, `UNIQUE` |
| `cart_items`   | `cart_id`   | `carts`       | `ON DELETE CASCADE`        |
| `cart_items`   | `product_id`| `products`    | `ON DELETE CASCADE`        |
| `orders`       | `user_id`   | `users`       | по умолчанию (`NO ACTION`) |
| `order_items`  | `order_id`  | `orders`      | `ON DELETE CASCADE`        |
| `order_items`  | `product_id`| `products`    | по умолчанию (`NO ACTION`) |

`cart_items` имеет первичный ключ `(cart_id, product_id)`:
повторное добавление товара увеличивает `quantity`. `orders` имеет уникальную
пару `(user_id, idempotency_key)` — повторный checkout с тем же ключом возвращает
уже существующий заказ вместо создания нового.

Обратите внимание: у пользователя статус берётся из `user_status`, а у продукта — из `status`.

## Таблицы

| Таблица       | Колонки                                                          |
|---------------|------------------------------------------------------------------|
| `status`      | `id` (PK), `name`, `code`                                        |
| `rules`       | `id` (PK), `name`, `code`                                        |
| `user_status` | `id` (PK), `name`, `code`                                        |
| `users`       | `id` (PK), `name`, `surname` (nullable), `female` (nullable), `email` (nullable, UNIQUE), `password` (nullable, bcrypt), `status_id` (FK), `rule_id` (FK), `created_at`, `updated_at` |
| `products`    | `id` (PK), `name`, `description` (nullable), `status_id` (FK), `price` (NUMERIC(10,2), ≥ 0), `stock` (INTEGER, ≥ 0) |
| `sessions`    | `id` (PK), `user_id` (FK, NOT NULL), `token_hash` (UNIQUE), `created_at`, `expires_at`, `revoked_at` (nullable) |
| `carts`       | `id` (PK), `user_id` (FK, UNIQUE), `created_at`, `updated_at` |
| `cart_items`  | `cart_id` (PK/FK), `product_id` (PK/FK), `quantity` (> 0)        |
| `orders`      | `id` (PK), `user_id` (FK), `status`, `total` (NUMERIC(10,2)), `idempotency_key` (VARCHAR(64), UNIQUE с `user_id`), `created_at` |
| `order_items` | `order_id` (PK/FK), `product_id` (PK/FK), `quantity` (> 0), `price` (NUMERIC(10,2), снимок) |

`id` — `SERIAL PRIMARY KEY`, то есть автоинкремент. `name` и `code` — `VARCHAR(255) NOT NULL`.
`female` — `VARCHAR(10)`, остальные строковые поля — `VARCHAR(255)`.
`created_at` и `expires_at` в `sessions` — `TIMESTAMPTZ NOT NULL`
(`created_at` по умолчанию `now()`), `revoked_at` — `TIMESTAMPTZ NULL`.

`users.created_at` и `users.updated_at` тоже `TIMESTAMPTZ NOT NULL` со
значением по умолчанию `now()` (миграции 000010 и 000011). При обновлении
пользователя бэкенд явно выставляет `updated_at = now()`; в ответе оба поля
приходят в формате RFC3339 и показываются в таблице пользователей панели.

`users.email` формально nullable (PostgreSQL разрешает несколько `NULL`), но
регистрация и сид админа всегда его заполняют; уникальность даёт ограничение
`users_email_key` из миграции 000009. `users.password` тоже nullable — у
пользователей без пароля вход невозможен.

В `sessions` хранится не сам токен, а его SHA-256-хеш (64 hex-символа).
Токен выдаётся клиенту при `/user/login` и живёт 24 часа; `revoked_at`
используется для выхода. Индексы `sessions_user_id_idx` и
`sessions_expires_at_idx` ускоряют поиск сессий пользователя и очистку
просроченных записей.

`products.price` и `products.stock` добавляются миграциями 000013 и 000015
со значением по умолчанию `0`. При оформлении заказа `stock` списывается в той же
транзакции, а `price` копируется в `order_items.price` — снимок цены на момент
покупки. `orders.total` считается БД как `SUM(quantity * price)` по позициям
заказа. `POST /cart/items` создаёт корзину при необходимости и увеличивает
`quantity` для уже существующей позиции; `DELETE /cart/items?product_id=` удаляет
одну позицию пользователя по `product_id`; `POST /order/checkout` очищает
`cart_items`, но саму корзину оставляет.

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
| 000010 | `000010_add_created_at_to_users.up.sql` | `users.created_at`             |
| 000011 | `000011_add_updated_at_to_users.up.sql` | `users.updated_at`             |
| 000012 | `000012_add_expires_at_index_sessions.up.sql` | индекс по `sessions.expires_at` |
| 000013 | `000013_add_price_to_products.up.sql`   | `products.price`                        |
| 000014 | `000014_create_carts.up.sql`            | `carts`, `cart_items`                   |
| 000015 | `000015_create_orders.up.sql`           | `products.stock`, `orders`, `order_items` |

Порядок важен: сначала таблицы без зависимостей, потом `users`, `products` и
`sessions`, которые ссылаются на них через `REFERENCES`; затем цена товаров,
корзина и заказы. Для миграций 000006–000015 есть `.down.sql`; у 000001–000005
down-файлов пока нет.

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
| `carts`        | `internal/cart`            |
| `orders`       | `internal/order`           |
| Пароли         | `internal/auth`            |
| Миграции, сид  | `internal/database`        |

В пакетах сущностей: `models.go` (структуры), `repository.go` (SQL),
`controller/` (HTTP-обработчики).

---

## Чек-лист

- [x] таблицы-справочники создаются раньше зависимых
- [x] связи описаны через `REFERENCES`
- [x] nullable-колонки в Go представлены указателями или `sql.Null*`
- [x] текстовое описание соответствует SQL-миграциям 000001–000015
- [ ] `docs/ERD.png` обновить под корзину и заказы (сейчас только базовые таблицы)
