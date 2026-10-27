package controller

import (
	"encoding/json"
	"fmt"
	"hello-go-backend/internal/product"
	"net/http"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Only POST allowed")
		return
	}
	var req product.CreateProductRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid JSON: %v", err)
		return
	}
	if req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "name is required")
		return
	}
	p, err := product.CreateProduct(h.db, req.Name, req.Description, req.StatusId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}
