package httpx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// ParseID достает id из query-параметра ?id=
func ParseID(r *http.Request) (int, error) {
	return ParseInt(r, "id")
}

func ParseInt(r *http.Request, key string) (int, error) {
	s := r.URL.Query().Get(key)
	if s == "" {
		return 0, fmt.Errorf("%s parameter is required", key)
	}
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return value, nil
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
