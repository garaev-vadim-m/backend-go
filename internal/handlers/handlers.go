package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"hello-go-backend/internal/database"
	"hello-go-backend/internal/models"
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
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error getting user: %v", err)
		return
	}

	// Возвращаем JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := database.GetAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Only post")
		return
	}

	// Парсим JSON из body
	var req models.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid request: %v", err)
		return
	}

	// Создаем пользователя
	user, err := database.CreateUser(&models.User{
		Name:     req.Name,
		Female:   req.Female,
		StatusID: req.StatusID,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error creating user: %v", err)
		return
	}

	// Возвращаем JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
