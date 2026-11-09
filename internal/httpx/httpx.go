package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// ParseID достает id из query-параметра ?id=
func ParseID(r *http.Request) (int, error) {
	s := r.URL.Query().Get("id")
	if s == "" {
		return 0, errors.New("id parameter is required")
	}
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid id: %w", err)
	}
	return id, nil
}

// WriteJSON пишет статус, заголовок и JSON-тело
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

type errorBody struct {
	Error string `json:"error"`
}

// WriteError пишет ошибку в JSON: {"error": "..."}
func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, errorBody{Error: err.Error()})
}
