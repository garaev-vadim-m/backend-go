package controller

import (
	"encoding/json"
	"fmt"
	"hello-go-backend/internal/user"
	"net/http"
	"strconv"
)

func (h *Handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Only PUT")
		return
	}

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
	var req user.UpdateUserRequest

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid request: %v", err)
		return
	}

	if req.Name == "" && req.StatusID == nil && req.RuleID == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "name, status, rule is required")
		return
	}

	user, err := user.UpdateUser(h.db, id, req.Name, req.Female, req.StatusID, req.RuleID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to update user: %v", err)
		return
	}

	// Отдаём JSON в том же формате, что и GET /users
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
