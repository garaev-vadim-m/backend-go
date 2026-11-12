# Ошибки в Go — errors.New, %w и errors.Is

## Ошибка — это значение

В Go ошибка — обычное возвращаемое значение, которое функция отдаёт последним:

```go
u, err := user.GetUser(db, id)
if err != nil {
    // обрабатываем ошибку
}
```

В JavaScript ошибку «бросают» и ловят через `try/catch`:

```javascript
try {
    const user = await getUser(id);
} catch (err) {
    console.log(err.message);
}
```

В Go `panic` оставляют для фатальных ситуаций (например, не удалось подключиться к БД при старте), а обычные ошибки возвращают и проверяют явно.

---

## errors.New — создать ошибку

```go
if req.Name == "" {
    httpx.WriteError(w, http.StatusBadRequest, errors.New("name is required"))
    return
}
```

Аналогия: `throw new Error("name is required")`.

---

## fmt.Errorf и %w — добавить контекст

```go
return nil, fmt.Errorf("user %d: %w", id, errs.ErrNotFound)
```

- `%w` — оборачивает (wrap) ошибку и сохраняет её внутри;
- `%v` — только подставляет текст, исходная ошибка теряется.

Обёрнутая ошибка остаётся «узнаваемой» через `errors.Is`.

---

## Ошибка-маркер в отдельном пакете

```go
package errs

import "errors"

// ErrNotFound означает, что сущность не найдена.
var ErrNotFound = errors.New("not found")
```

Репозитории добавляют контекст (какой именно id не найден), а сравнивают контроллеры.

---

## Цепочка обработки

```
repository
  sql.ErrNoRows
      │
      │ fmt.Errorf("user %d: %w", id, errs.ErrNotFound)
      ▼
controller
  errors.Is(err, errs.ErrNotFound) → true
      │
      ▼
HTTP 404 Not Found
```

В коде контроллера:

```go
u, err := user.GetUser(h.db, id)
if err != nil {
    if errors.Is(err, errs.ErrNotFound) {
        httpx.WriteError(w, http.StatusNotFound, err)
        return
    }
    log.Printf("user get: %v", err)
    httpx.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
    return
}
```

---

## Что использовать

| Функция                     | Назначение                                |
|-----------------------------|-------------------------------------------|
| `errors.New("...")`         | создать простую ошибку                    |
| `fmt.Errorf("... %w", err)` | обернуть ошибку с контекстом              |
| `errors.Is(err, target)`    | проверить, есть ли `target` в цепочке     |
| `errors.As(err, &target)`   | привести ошибку к конкретному типу        |

---

## Чек-лист

- [ ] функция возвращает `error` последним значением
- [ ] ошибку проверяем через `if err != nil`
- [ ] контекст добавляем через `fmt.Errorf("...: %w", err)`
- [ ] тип проверяем через `errors.Is(err, errs.ErrNotFound)`
- [ ] клиенту отдаём `{"error": "..."}`, детали пишем в `log`
