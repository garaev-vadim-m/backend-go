package handlers

import (
	"fmt"
	"net/http"

	"hello-go-backend/internal/database"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Go Backend!")
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if err := database.DB.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Database error: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}
