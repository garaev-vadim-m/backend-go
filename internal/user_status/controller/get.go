package controller

import (
	"errors"
	"hello-go-backend/internal/errs"
	"hello-go-backend/internal/httpx"
	userstatus "hello-go-backend/internal/user_status"
	"log"
	"net/http"
)

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ParseID(r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err)
		return
	}

	st, err := userstatus.GetStatus(h.db, id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			httpx.WriteError(w, http.StatusNotFound, err)
			return
		}
		log.Printf("user status get: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}

	httpx.WriteJSON(w, http.StatusOK, st)
}
