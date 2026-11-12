# Схема базы данных

## ERD

Визуальная схема лежит в `docs/ERD.png`. Текстовый вариант:

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
                │ female          │
                │ status_id  FK   │
                │ rule_id    FK   │
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

| Откуда         | Поле        | Куда          |
|----------------|-------------|---------------|
| `users`        | `status_id` | `user_status` |
| `users`        | `rule_id`   | `rules`       |
| `products`     | `status_id` | `status`      |

Обратите внимание: у пользователя статус берётся из `user_status`, а у продукта — из `status`.

## Таблицы

| Таблица       | Колонки                                                          |
|---------------|------------------------------------------------------------------|
| `status`      | `id` (PK), `name`, `code`                                        |
| `rules`       | `id` (PK), `name`, `code`                                        |
| `user_status` | `id` (PK), `name`, `code`                                        |
| `users`       | `id` (PK), `name`, `female` (nullable), `status_id` (FK), `rule_id` (FK) |
| `products`    | `id` (PK), `name`, `description` (nullable), `status_id` (FK)    |

`id` — `SERIAL PRIMARY KEY`, то есть автоинкремент. `name` и `code` — `VARCHAR(255) NOT NULL`.

## Миграции

Таблицы создаются в `internal/database/migrations.go` при старте приложения:

```go
createStatusTable := `
CREATE TABLE IF NOT EXISTS status (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(255) NOT NULL
);
`
```

Все запросы собраны в срез:

```go
tables := []struct {
    name string
    sql  string
}{
    {"status", createStatusTable},
    {"rules", createRulesTable},
    {"user_status", createUserStatusTable},
    {"users", createUsersTable},
    {"products", createProductsTable},
}
```

Порядок важен: сначала таблицы без зависимостей, потом `users` и `products`, которые ссылаются на них через `REFERENCES`.

## Схема и пакеты

| Таблица       | Go-пакет                  |
|---------------|---------------------------|
| `status`      | `internal/status`         |
| `rules`       | `internal/rules`          |
| `user_status` | `internal/user_status`    |
| `users`       | `internal/user`           |
| `products`    | `internal/product`        |

В каждом пакете: `models.go` (структуры), `repository.go` (SQL), `controller/` (HTTP-обработчики).

---

## Чек-лист

- [ ] таблицы-справочники создаются раньше зависимых
- [ ] связи описаны через `REFERENCES`
- [ ] nullable-колонки в Go представлены указателями или `sql.Null*`
- [ ] схема совпадает с `docs/ERD.png`
