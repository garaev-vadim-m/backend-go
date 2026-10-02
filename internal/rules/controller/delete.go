package controller

import (
	"fmt"
	"hello-go-backend/internal/rules"
	"net/http"
	"strconv"
)

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Only DELETE allowed")
		return
	}

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

	err = rules.DeleteRules(h.db, id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
