package middleware

import (
	"errors"
	"hello-go-backend/internal/httpx"
	"log"
	"net/http"
)

func (a *Auth) RequireRole(code string, next http.Handler) http.Handler {
	return a.Require(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := UserIDFromContext(r.Context())

		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, errors.New("not logged in"))
			return
		}

		var role string
		err := a.db.QueryRow(`
			SELECT r.code FROM users u
			LEFT JOIN rules r ON r.id = u.rule_id
			WHERE u.id = $1`, userID).Scan(&role)

		if err != nil {
			log.Printf("role check: %v", err)
			httpx.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
			return
		}

		if role != code {
			httpx.WriteError(w, http.StatusForbidden, errors.New("forbidden"))
		}

		next.ServeHTTP(w, r)
	}))
}
