package handlers

import (
	"encoding/json"
	"fmt"
	"hello-go-backend/internal/database"
	"net/http"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	statuses, err := database.GetAllStatus()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statuses)
}
