package controller

import (
	"encoding/json"
	"fmt"
	"hello-go-backend/internal/rules"
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

	var req rules.UpdateRuleRequest

	if err = json.NewDecoder(r.Body).Decode(&req); req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "name is required")
		return
	}

	st, err := rules.UpdateRule(h.db, id, req.Name, req.Code)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(st)
}
