# database/sql — работа с PostgreSQL

## Подключение

```go
import (
    "database/sql"

    _ "github.com/lib/pq" // драйвер регистрируется сам
)

db, err := sql.Open("postgres", connStr)
if err != nil {
    return err
}
defer db.Close()
```

`_` перед импортом означает: пакет нужен только ради побочного эффекта — регистрации драйвера `postgres` в `database/sql`. Напрямую его функции мы не вызываем.

## Query, QueryRow, Exec

| Метод       | Когда использовать                 | Результат    |
|-------------|------------------------------------|--------------|
| `Query`     | SELECT, много строк                | `*sql.Rows`  |
| `QueryRow`  | SELECT/INSERT/UPDATE, одна строка  | `*sql.Row`   |
| `Exec`      | DELETE/UPDATE без возврата данных  | `sql.Result` |

## Плейсхолдеры и SQL-инъекции

В PostgreSQL параметры пишутся через `$1`, `$2`, ...:

```go
query := `SELECT id, name, code FROM status WHERE id = $1`
err := db.QueryRow(query, id).Scan(&st.ID, &st.Name, &st.Code)
```

Драйвер сам экранирует значения — это защита от SQL-инъекций (как параметризованные запросы в `pg` для Node.js).

```javascript
// JS (pg)
await client.query('SELECT id, name, code FROM status WHERE id = $1', [id]);
```

## Чтение нескольких строк

```go
rows, err := db.Query(query)
if err != nil {
    return nil, err
}
defer rows.Close() // обязательно закрываем

var statuses []*Status
for rows.Next() {
    st := &Status{}
    if err := rows.Scan(&st.ID, &st.Name, &st.Code); err != nil {
        return nil, err
    }
    statuses = append(statuses, st)
}

if err := rows.Err(); err != nil { // ошибка во время итерации
    return nil, err
}
```

## RowsAffected — проверить, что строка была

```go
result, err := db.Exec(`DELETE FROM users WHERE id = $1`, id)
if err != nil {
    return err
}

n, err := result.RowsAffected()
if err != nil {
    return err
}
if n == 0 {
    return fmt.Errorf("user %d: %w", id, errs.ErrNotFound)
}
```

## sql.ErrNoRows — «ничего не нашли»

```go
err := db.QueryRow(query, id).Scan(&st.ID, &st.Name, &st.Code)
if err == sql.ErrNoRows {
    return nil, fmt.Errorf("status %d: %w", id, errs.ErrNotFound)
}
```

`QueryRow` не возвращает пустую строку — при отсутствии данных `Scan` отдаёт `sql.ErrNoRows`.

## NULL-колонки: sql.Null*

Для JOIN-нутых полей статуса и роли используем `sql.NullInt64` и `sql.NullString`:

```go
var statusID sql.NullInt64
var statusName, statusCode sql.NullString

err := s.Scan(&statusID, &statusName, &statusCode)
if err != nil {
    return nil, err
}

if statusID.Valid {
    p.Status = &status.Status{
        ID:   int(statusID.Int64),
        Name: statusName.String,
        Code: statusCode.String,
    }
}
```

Поле `Valid` (было ли значение) и `.Int64` / `.String` (само значение) позволяют отличить `NULL` от нуля и пустой строки.

## RETURNING + CTE + JOIN

`RETURNING` не умеет `JOIN`, поэтому в проекте INSERT/UPDATE заворачивают в CTE:

```sql
WITH ins AS (
    INSERT INTO products (name, description, status_id)
    VALUES ($1, $2, $3)
    RETURNING id, name, description, status_id
)
SELECT
    p.id, p.name, p.description,
    st.id, st.name, st.code
FROM ins p
LEFT JOIN status st ON st.id = p.status_id
```

Так ответ сразу содержит развёрнутый статус, а не только `status_id`.

### LEFT JOIN и nullable

`LEFT JOIN` возвращает строку даже без связанной записи — колонки статуса будут `NULL`. Поэтому связанный объект в модели — указатель (`*status.Status`), а в JSON получается `"status": null`.

---

## Чек-лист

- [ ] драйвер подключён через `_ "github.com/lib/pq"`
- [ ] параметры передаём через `$1`, `$2`, а не склейкой строк
- [ ] `defer rows.Close()` после `Query`
- [ ] после цикла проверяем `rows.Err()`
- [ ] `sql.ErrNoRows` превращаем в `errs.ErrNotFound`
- [ ] для NULL-колонок используем `sql.Null*` и проверяем `.Valid`
- [ ] INSERT/UPDATE с JOIN делаем через CTE + `RETURNING`
