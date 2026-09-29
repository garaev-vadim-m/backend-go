package controller

import (
	"encoding/json"
	"fmt"
	"hello-go-backend/internal/status"
	"net/http"
)

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	statuses, err := status.GetAllStatus(h.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statuses)
}
