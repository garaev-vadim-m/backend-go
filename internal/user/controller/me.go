package controller

import (
	"errors"
	"hello-go-backend/internal/httpx"
	"hello-go-backend/internal/middleware"
	"hello-go-backend/internal/user"
	"log"
	"net/http"
)

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	me, ok := middleware.IdentityFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, errors.New("not logged in"))
		return
	}

	u, err := user.GetUser(h.db, me.ID)
	if err != nil {
		log.Printf("user me get: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}

	httpx.WriteJSON(w, http.StatusOK, u)
}
