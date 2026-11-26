package controller

import (
	"encoding/json"
	"errors"
	"hello-go-backend/internal/cart"
	"hello-go-backend/internal/httpx"
	"hello-go-backend/internal/middleware"
	"log"
	"net/http"

	"github.com/lib/pq"
)

type addItemRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	me, ok := middleware.IdentityFromContext(r.Context())

	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, errors.New("not logged in"))
		return
	}

	var req addItemRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, errors.New("invalid JSON"))
		return
	}

	if req.ProductID <= 0 || req.Quantity <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, errors.New("product_id and quantity must be positive"))
		return
	}

	if err := cart.AddItem(h.db, me.ID, req.ProductID, req.Quantity); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			httpx.WriteError(w, http.StatusNotFound, errors.New("product not found"))
			return
		}
		log.Printf("cart add item: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}
}
