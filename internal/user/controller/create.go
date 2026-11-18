package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"hello-go-backend/internal/auth"
	"hello-go-backend/internal/httpx"
	"hello-go-backend/internal/user"
	"log"
	"net/http"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req user.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid JSON: %w", err))
		return
	}

	if req.Name == "" {
		httpx.WriteError(w, http.StatusBadRequest, errors.New("name is required"))
		return
	}
	if req.Password == nil || *req.Password == "" {
		httpx.WriteError(w, http.StatusBadRequest, errors.New("password is required"))
		return
	}

	passwordHash, err := auth.HashOptional(req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrWeakPassword) {
			httpx.WriteError(w, http.StatusBadRequest, err)
			return
		}
		log.Printf("user create hash: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}

	u, err := user.CreateUser(h.db, req.Name, req.Surname, req.Female, req.Email, passwordHash, req.StatusID, req.RuleID)
	if err != nil {
		log.Printf("user create: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, u)
}
