package controller

import (
	"encoding/json"
	"fmt"
	"hello-go-backend/internal/status"
	"net/http"
	"strconv"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "id parameter is required")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid id: %v", err)
		return
	}
	var req status.UpdateStatusRequest

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid JSON: %v", err)
		return
	}
	if req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "name is required")
		return
	}
	st, err := status.UpdateStatus(h.db, id, req.Name, req.Code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(st)
}
