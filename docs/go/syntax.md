# Go — базовый синтаксис

## defer

`defer` выполняет функцию в конце текущей функции (вне зависимости от результата).

```go
func main() {
    db, _ := sql.Open("postgres", connStr)
    defer db.Close() // выполнится при выходе из main

    // ... работа с БД ...
}
```

Полезно для освобождения ресурсов: закрытие файлов, соединений, mutex'ов.

## package main

```go
package main
```

Точка входа в Go-программу. Файл с `package main` и функцией `main()` — это исполняемый бинарник.

## := — объявление переменной

```go
name := "John"    // строка
age := 30         // int
```

Короткая запись объявления с автоматическим выводом типа.

## func — функция

```go
func add(a int, b int) int {
    return a + b
}
```

## for range N — цикл без счётчика (Go 1.22)

```go
for i := range 5 {
    // i = 0, 1, 2, 3, 4
}
```

Аналогия в JavaScript: `for (let i = 0; i < 5; i++) {}`. Объявлять счётчик и писать инкремент не нужно.

## struct — структура

```go
type Status struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Code string `json:"code"`
}
```

Структура — набор именованных полей с фиксированными типами (как объект в JS, но состав полей задан заранее). Теги `json:"..."` описаны в `json.md`.

## Методы и указательный получатель

```go
type Handler struct {
    db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
    return &Handler{db: db}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
    // h — аналог this в JS
}
```

- `(h *Handler)` — получатель метода: внутри доступен `h`, как `this` в объекте;
- `*Handler` — указатель, чтобы все методы работали с одним и тем же `Handler` и его `db`;
- `NewHandler` — конструктор: в Go нет классов, вместо них структура + функции.

Аналогия в JavaScript:

```javascript
class Handler {
    constructor(db) { this.db = db; }
    get(req, res) { /* this.db */ }
}
```

## Множественные возвращаемые значения

```go
u, err := user.GetUser(db, id)
if err != nil {
    // ...
}
```

Функция может вернуть несколько значений. Частая пара — результат + ошибка.

Аналогия в JavaScript: `const [user, err] = await getUser(id)`.

## panic vs error

```go
if err := database.InitDB(); err != nil {
    panic(err) // фатально: без БД сервер работать не может
}
```

- `error` — обычная ошибка, её возвращают и обрабатывают;
- `panic` — «всё сломалось», программу лучше остановить (например, на старте).
