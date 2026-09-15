Окей, пишем функцию для получения пользователя.

---

## 1. Создаем internal/database/user.go

```go
package database

import (
	"database/sql"
	"fmt"

	"hello-go-backend/internal/models"
)

// GetUser получает пользователя по ID из БД
func GetUser(id int) (*models.User, error) {
	// SQL запрос
	query := `SELECT id, name, female, status_id FROM users WHERE id = $1`

	// Создаем переменную для результата
	user := &models.User{}

	// Выполняем запрос
	err := DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Female,
		&user.StatusID,
	)

	// Проверяем ошибки
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, err
	}

	return user, nil
}
```

**Что здесь:**
- `$1` — плейсхолдер для параметра (защита от SQL injection)
- `QueryRow()` — выполняет запрос, возвращает одну строку
- `Scan()` — заполняет переменные из результата
- `sql.ErrNoRows` — специальная ошибка если пользователя нет

---

## 2. Создаем обработчик в handlers

Обновляешь `internal/handlers/handlers.go`:

```go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"hello-go-backend/internal/database"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Go Backend!")
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if err := database.DB.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Database error: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

// Новый обработчик
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из URL параметра
	// URL: /user/1 -> query параметр "id"
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "id parameter is required")
		return
	}

	// Преобразуем строку в число
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid id: %v", err)
		return
	}

	// Получаем пользователя из БД
	user, err := database.GetUser(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	// Возвращаем JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
```

**Что здесь:**
- `r.URL.Query().Get("id")` — получить параметр из URL (?id=1)
- `strconv.Atoi()` — преобразовать строку в int
- `json.NewEncoder()` — кодировать структуру в JSON

---

## 3. Регистрируем маршрут в main.go

```go
package main

import (
	"fmt"
	"net/http"

	"hello-go-backend/internal/database"
	"hello-go-backend/internal/handlers"
)

func main() {
	if err := database.InitDB(); err != nil {
		panic(err)
	}
	defer database.CloseDB()

	if err := database.RunMigrations(); err != nil {
		panic(err)
	}

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/user", handlers.GetUserHandler)  // ← новый маршрут

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
```

---

## 4. Запускаем и проверяем

```bash
docker-compose down
docker-compose up --build
```

---

## 5. Тестируем

Сначала добавим тестовых данных:

```bash
docker-compose exec postgres psql -U postgres -d myapp
```

Внутри psql:
```sql
INSERT INTO users (name, female, status_id) VALUES ('John', 'Male', 1);
INSERT INTO users (name, female, status_id) VALUES ('Jane', 'Female', NULL);
```

Потом:
```bash
curl "http://localhost:8080/user?id=1"
```

Ответ должен быть:
```json
{"id":1,"name":"John","female":"Male","status_id":1}
```

Несуществующий пользователь:
```bash
curl "http://localhost:8080/user?id=999"
```

Ответ:
```
Error: user with id 999 not found
```

---
