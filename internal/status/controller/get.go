package controller

import (
	"encoding/json"
	"fmt"
	"hello-go-backend/internal/status"
	"net/http"
	"strconv"
)

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
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
	st, err := status.GetStatus(h.db, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(st)
}
