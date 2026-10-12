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

	if req.Name == "" && req.RuleID == 0 && req.StatusID == 0 {
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

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "user updated: %v", user)
}
