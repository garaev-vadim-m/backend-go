package controller

import (
	"encoding/json"
	"fmt"
	"hello-go-backend/internal/user"
	"net/http"
)

func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Only post")
		return
	}

	// Парсим JSON из body
	var req user.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid request: %v", err)
		return
	}

	// Создаем пользователя
	user, err := user.CreateUser(h.db, &user.User{
		Name:     req.Name,
		Female:   req.Female,
		StatusID: req.StatusID,
		RuleID:   req.RuleID,
	})

	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "error creating user: %v", err)
		return
	}

	// Возвращаем JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
