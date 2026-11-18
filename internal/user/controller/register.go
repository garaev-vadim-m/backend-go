package controller

import (
	"encoding/json"
	"errors"
	"hello-go-backend/internal/auth"
	"hello-go-backend/internal/httpx"
	"hello-go-backend/internal/rules"
	"hello-go-backend/internal/user"
	"log"
	"net/http"

	"github.com/lib/pq"
)

// registerRequest намеренно не содержит rule_id и status_id:
// регистрация не должна позволять клиенту выбрать свою роль.
type registerRequest struct {
	Name     string  `json:"name"`
	Surname  *string `json:"surname"`
	Female   *string `json:"female"`
	Email    string  `json:"email"`
	Password *string `json:"password"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, errors.New("name and email are required"))
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
		log.Printf("user register hash: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}

	userRuleID, err := rules.GetIDByCode(h.db, "user")
	if err != nil {
		log.Printf("user register default rule: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}

	email := req.Email
	u, err := user.CreateUser(h.db, req.Name, req.Surname, req.Female, &email, passwordHash, nil, &userRuleID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			httpx.WriteError(w, http.StatusConflict, errors.New("email already registered"))
			return
		}
		log.Printf("user register create: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, u)
}
