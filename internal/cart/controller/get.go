package controller

import (
	"errors"
	"hello-go-backend/internal/cart"
	"hello-go-backend/internal/httpx"
	"hello-go-backend/internal/middleware"
	"log"
	"net/http"
)

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	me, ok := middleware.IdentityFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, errors.New("not logged in"))
		return
	}

	c, err := cart.GetCart(r.Context(), h.db, me.ID)
	if err != nil {
		log.Printf("cart get: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}

	httpx.WriteJSON(w, http.StatusOK, c)
}
