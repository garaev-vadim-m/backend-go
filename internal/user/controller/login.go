package controller

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"hello-go-backend/internal/auth"
	"hello-go-backend/internal/errs"
	"hello-go-backend/internal/httpx"
	"hello-go-backend/internal/session"
	"hello-go-backend/internal/user"
	"log"
	"net/http"
	"time"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, errors.New("invalid JSON"))
		return
	}

	if req.Email == "" || req.Password == "" {
		httpx.WriteError(w, http.StatusBadRequest, errors.New("email and password are required"))
		return
	}

	id, stored, err := user.GetUserPasswordByEmail(h.db, req.Email)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			// Одинаковый ответ для «нет пользователя» и «неверный пароль»,
			// чтобы не давать перебирать почты.
			httpx.WriteError(w, http.StatusUnauthorized, errors.New("invalid credentials"))
			return
		}
		log.Printf("user login: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}

	if !auth.CheckPassword(stored, req.Password) && subtle.ConstantTimeCompare([]byte(stored), []byte(req.Password)) != 1 {
		httpx.WriteError(w, http.StatusUnauthorized, errors.New("invalid credentials"))
		return
	}

	token, err := session.Create(h.db, id, 24*time.Hour)
	if err != nil {
		log.Printf("user login get: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}

	session.SetCookie(w, token, 24*time.Hour, h.cookieSecure)
	httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
