package controller

import (
	"errors"
	"log"
	"net/http"

	"hello-go-backend/internal/errs"
	"hello-go-backend/internal/httpx"
	"hello-go-backend/internal/middleware"
	"hello-go-backend/internal/user"
)

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ParseID(r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err)
		return
	}

	me, ok := middleware.IdentityFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, errors.New("not logged in"))
		return
	}
	if me.Role != "admin" && me.ID != id {
		//403 удобно для отладки
		// httpx.WriteError(w, http.StatusForbidden, errors.New("forbidden"))
		// 404 чтобы пользователь без роли admin не мог понять что пользователь вообще существует
		httpx.WriteError(w, http.StatusNotFound, errors.New("forbidden"))
		return
	}

	u, err := user.GetUser(h.db, id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			httpx.WriteError(w, http.StatusNotFound, err)
			return
		}
		log.Printf("user get: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}

	httpx.WriteJSON(w, http.StatusOK, u)
}
