package controller

import (
	"encoding/json"
	"fmt"
	"hello-go-backend/internal/rules"
	"net/http"
)

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	rules, err := rules.GetAllRules(h.db)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rules)
}
