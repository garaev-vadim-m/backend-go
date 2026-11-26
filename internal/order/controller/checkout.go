package controller

import (
	"errors"
	"hello-go-backend/internal/errs"
	"hello-go-backend/internal/httpx"
	"hello-go-backend/internal/middleware"
	"hello-go-backend/internal/order"
	"log"
	"net/http"
)

func (h *Handler) Checkout(w http.ResponseWriter, r *http.Request) {
	me, ok := middleware.IdentityFromContext(r.Context())

	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, errors.New("not logged in"))
		return
	}

	key := r.Header.Get("Idempotency-Key")

	if key == "" || len(key) > 64 {
		httpx.WriteError(w, http.StatusBadRequest, errors.New("idempotency-key header is required (max 64 chars)"))
		return
	}

	o, created, err := order.Checkout(r.Context(), h.db, me.ID, key)

	if err != nil {
		switch {
		case errors.Is(err, errs.ErrEmptyCart):
			httpx.WriteError(w, http.StatusBadRequest, err)
		case errors.Is(err, errs.ErrOutOfStock):
			httpx.WriteError(w, http.StatusConflict, err)
		default:
			log.Printf("order checkout: %v", err)
			httpx.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
		}
		return
	}

	status := http.StatusOK
	if created {
		status = http.StatusCreated
	}
	httpx.WriteJSON(w, status, o)
}
