package controller

import (
	"fmt"
	"hello-go-backend/internal/product"
	"net/http"
	"strconv"
)

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "id is required")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid id: %v", err)
		return
	}

	if err = product.DeleteProduct(h.db, id); err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
