package controller

import (
	"errors"
	"log"
	"net/http"

	"hello-go-backend/internal/httpx"
	"hello-go-backend/internal/session"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(session.CookieName)
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, errors.New("not logged in"))
		return
	}

	if err := session.Revoke(h.db, c.Value); err != nil {
		log.Printf("user logout: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}

	session.ClearCookie(w, h.cookieSecure)
	w.WriteHeader(http.StatusNoContent)
}
