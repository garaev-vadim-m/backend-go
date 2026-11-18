package controller

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"hello-go-backend/internal/httpx"
	"hello-go-backend/internal/session"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !ok || token == "" {
		httpx.WriteError(w, http.StatusUnauthorized, errors.New("missing token"))
		return
	}

	if err := session.Revoke(h.db, token); err != nil {
		log.Printf("user logout: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
