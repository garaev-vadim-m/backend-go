package handlers

import (
	"encoding/json"
	"fmt"
	"hello-go-backend/internal/database"
	"net/http"
	"strconv"
)

func StatusesHandler(w http.ResponseWriter, r *http.Request) {
	statuses, err := database.GetAllStatus()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statuses)
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "id is required")
		return
	}

	id, err := strconv.Atoi(idStr)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	status, err := database.GetStatus(id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
