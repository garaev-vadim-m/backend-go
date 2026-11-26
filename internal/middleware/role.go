package middleware

import (
	"errors"
	"hello-go-backend/internal/httpx"
	"net/http"
)

func (a *Auth) RequireRole(code string, next http.Handler) http.Handler {
	return a.Require(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		me, ok := IdentityFromContext(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, errors.New("not logged in"))
			return
		}
		if me.Role != code {
			httpx.WriteError(w, http.StatusForbidden, errors.New("forbidden"))
			return
		}
		next.ServeHTTP(w, r)
	}))
}
