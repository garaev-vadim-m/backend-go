package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"hello-go-backend/internal/errs"
	"hello-go-backend/internal/httpx"
	"hello-go-backend/internal/status"
	"log"
	"net/http"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ParseID(r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var req status.UpdateStatusRequest

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid JSON: %w", err))
		return
	}
	if req.Name == "" {
		httpx.WriteError(w, http.StatusBadRequest, errors.New("name is required"))
		return
	}
	st, err := status.UpdateStatus(h.db, id, req.Name, req.Code)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			httpx.WriteError(w, http.StatusNotFound, err)
			return
		}
		log.Printf("status update: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}

	httpx.WriteJSON(w, http.StatusOK, st)
}
