package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"hello-go-backend/internal/httpx"
	"hello-go-backend/internal/rules"
	"log"
	"net/http"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req rules.CreateRuleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid JSON: %w", err))
		return
	}

	if req.Name == "" {
		httpx.WriteError(w, http.StatusBadRequest, errors.New("name is required"))
		return
	}

	st, err := rules.CreateRule(h.db, req.Name, req.Code)
	if err != nil {
		log.Printf("rule create: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, st)
}
