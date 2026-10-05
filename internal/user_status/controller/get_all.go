package controller

import (
	"encoding/json"
	"fmt"
	userstatus "hello-go-backend/internal/user_status"
	"net/http"
)

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	user_statuses, err := userstatus.GetAllStatus(h.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user_statuses)
}
