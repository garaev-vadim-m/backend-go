package controller

import (
	"errors"
	"hello-go-backend/internal/cart"
	"hello-go-backend/internal/httpx"
	"hello-go-backend/internal/middleware"
	"log"
	"net/http"
)

func (h *Handler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	me, ok := middleware.IdentityFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, errors.New("not logged in"))
		return
	}

	productID, err := httpx.ParseInt(r, "product_id")
	if err != nil || productID <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, errors.New("positive product_id is required"))
		return
	}

	removed, err := cart.RemoveItem(h.db, me.ID, productID)
	if err != nil {
		log.Printf("cart remove item: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}

	if !removed {
		httpx.WriteError(w, http.StatusNotFound, errors.New("item not found in cart"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
