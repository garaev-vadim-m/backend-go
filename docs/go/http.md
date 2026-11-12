# net/http — сервер, роутинг и middleware

## Запуск сервера

```go
fmt.Println("Server running on :8080")
if err := http.ListenAndServe(":8080", corsMiddleware(mux)); err != nil {
    panic(err)
}
```

`http.ListenAndServe` принимает адрес и обработчик (router). Второй аргумент — это `http.Handler`, у которого есть метод `ServeHTTP`.

## ServeMux и роутинг (Go 1.22)

```go
mux := http.NewServeMux()

mux.HandleFunc("GET /user", userHandler.Get)
mux.HandleFunc("GET /users", userHandler.GetAll)
mux.HandleFunc("POST /user/create", userHandler.Create)
mux.HandleFunc("DELETE /user/delete", userHandler.Delete)
mux.HandleFunc("PUT /user/update", userHandler.Update)
```

Начиная с Go 1.22 метод пишется прямо в шаблоне: `"МЕТОД /путь"`.

Аналогия в JavaScript (Express):

```javascript
app.get('/user', userHandler.get);
app.post('/user/create', userHandler.create);
```

## ResponseWriter и Request

```go
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
    // w — куда пишем ответ
    // r — входящий запрос (URL, заголовки, тело)
}
```

| Go                    | JS (Express)  |
|-----------------------|---------------|
| `http.ResponseWriter` | `res`         |
| `*http.Request`       | `req`         |

## Handler-структура и конструктор

```go
type Handler struct {
    db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
    return &Handler{db: db}
}
```

В `main.go` создаётся один обработчик на сущность:

```go
userHandler := controllerUser.NewHandler(database.DB)
```

`(h *Handler)` в методах — это аналог `this` в JS: все методы работают с одним и тем же `h.db`.

## http.HandlerFunc

`http.HandlerFunc` превращает обычную функцию в `http.Handler`:

```go
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // ...
})
```

Именно так устроен `corsMiddleware`.

## Middleware

Middleware — обёртка над обработчиком:

```go
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // что-то до
        next.ServeHTTP(w, r)
        // что-то после
    })
}
```

Схема вызова:

```
http.ListenAndServe
        │
        ▼
corsMiddleware (заголовки + OPTIONS)
        │
        ▼
     ServeMux (роутинг по методу и пути)
        │
        ▼
   Handler.Get / Create / ...
        │
        ▼
   repository → PostgreSQL
```

## CORS

Браузер не пустит запрос с другого origin (`file://`, другой порт), если сервер не разрешил его явно:

```go
w.Header().Set("Access-Control-Allow-Origin", origin)
w.Header().Set("Vary", "Origin")
w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")

if r.Method == http.MethodOptions {
    w.WriteHeader(http.StatusNoContent) // preflight
    return
}
```

- `Access-Control-Allow-Origin` — с какого origin разрешены запросы;
- `Access-Control-Allow-Methods` — какие методы разрешены;
- `Access-Control-Allow-Headers` — какие заголовки можно присылать;
- `Vary: Origin` — ответ зависит от заголовка `Origin` (нужно для кешей);
- `OPTIONS` — preflight-запрос, браузер шлёт его перед `PUT`/`DELETE`.

## Коды ответов в проекте

| Код | Когда                                   |
|-----|-----------------------------------------|
| 200 | успешный `GET` / `PUT`                  |
| 201 | создание (`POST .../create`)            |
| 204 | удаление (`DELETE`)                     |
| 400 | битый JSON или отсутствует `id`/`name`  |
| 404 | сущность не найдена (`errs.ErrNotFound`)|
| 500 | внутренняя ошибка сервера               |

## Хелперы httpx

```go
id, err := httpx.ParseID(r)                        // ?id= из query
httpx.WriteJSON(w, http.StatusOK, u)               // JSON-ответ
httpx.WriteError(w, http.StatusNotFound, err)      // {"error": "..."}
```

`ParseID` достаёт `id` из query-строки (`GET /user?id=1`), проверяет, что он есть, и конвертирует в `int` через `strconv.Atoi`.

---

## Чек-лист

- [ ] маршруты зарегистрированы как `"МЕТОД /путь"`
- [ ] обработчики — методы на `*Handler`, создаются через `NewHandler`
- [ ] `Content-Type: application/json` выставляется в `httpx.WriteJSON`
- [ ] middleware вызывается как `corsMiddleware(mux)`
- [ ] preflight `OPTIONS` отвечает `204`
- [ ] коды ответов соответствуют таблице выше
