# encoding/json — теги, декодирование и ответы

## Теги структуры

```go
type Status struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Code string `json:"code"`
}
```

Тег `json:"..."` задаёт имя поля в JSON. Без тега Go взял бы имя поля как есть (`ID`, `Name`), а тег делает его `id`, `name`.

## Декодирование запроса

```go
var req product.CreateProductRequest

if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    httpx.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid JSON: %w", err))
    return
}
```

Аналогия в JavaScript:

```javascript
const req = await request.json(); // Express: req.body
```

`Decode` принимает указатель `&req`, чтобы заполнить структуру.

## Кодирование ответа

```go
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(p)
```

Аналогия: `res.status(201).json(p)`.

В проекте это инкапсулировано в `httpx.WriteJSON`:

```go
httpx.WriteJSON(w, http.StatusCreated, p)
httpx.WriteError(w, http.StatusBadRequest, errors.New("name is required"))
```

`WriteError` отдаёт `{"error": "..."}`.

## Указатели = необязательные поля

```go
type CreateUserRequest struct {
    Name     string  `json:"name"`
    Female   *string `json:"female"`
    StatusID *int    `json:"status_id"`
    RuleID   *int    `json:"rule_id"`
}
```

- `string` — поле обязательное, отсутствие даст `""`;
- `*string` / `*int` — если в JSON `null` или поля нет, будет `nil`;
- `nil` при записи в БД превращается в `NULL`.

## Запрос и ответ — разные структуры

| Структура           | Назначение             | Поля                                     |
|---------------------|------------------------|------------------------------------------|
| `CreateUserRequest` | разобрать тело `POST`  | `name`, `female`, `status_id`, `rule_id` |
| `UpdateUserRequest` | разобрать тело `PUT`   | `name`, `female`, `status_id`, `rule_id` |
| `User`              | отдать клиенту         | `id`, `name`, `female`, `status`, `rule` |

`status_id` / `rule_id` передаются только на запись. В ответе вместо них — развёрнутые объекты `status` и `rule`.

## Типы и их JSON

| Go                | JSON                                       |
|-------------------|--------------------------------------------|
| `int`, `string`   | число, строка                              |
| `*string` = nil   | `null`                                     |
| `*int` = nil      | `null`                                     |
| `*Status` = nil   | `null`                                     |
| `*Status` != nil  | `{ "id": ..., "name": ..., "code": ... }`  |
| `[]*User`         | массив объектов                            |

---

## Чек-лист

- [ ] у полей модели есть теги `json:"..."`
- [ ] декодируем через `json.NewDecoder(r.Body).Decode(&req)`
- [ ] отвечаем через `httpx.WriteJSON` / `httpx.WriteError`
- [ ] необязательные поля — указатели (`*string`, `*int`)
- [ ] для записи и чтения используются разные структуры
