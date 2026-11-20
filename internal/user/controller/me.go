package controller

import (
	"errors"
	"hello-go-backend/internal/errs"
	"hello-go-backend/internal/httpx"
	"hello-go-backend/internal/session"
	"hello-go-backend/internal/user"
	"log"
	"net/http"
)

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(session.CookieName)
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, errors.New("not logged in"))
		return
	}

	userID, err := session.Validate(h.db, c.Value)
	if errors.Is(err, errs.ErrUnauthorized) {
		httpx.WriteError(w, http.StatusUnauthorized, errors.New("not logged in"))
		return
	}
	if err != nil {
		log.Printf("user me validate: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}

	u, err := user.GetUser(h.db, userID)
	if err != nil {
		log.Printf("user me get: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}

	httpx.WriteJSON(w, http.StatusOK, u)
}
