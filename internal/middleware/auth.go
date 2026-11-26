package middleware

import (
	"context"
	"database/sql"
	"errors"
	"hello-go-backend/internal/errs"
	"hello-go-backend/internal/httpx"
	"hello-go-backend/internal/session"
	"log"
	"net/http"
	"strings"
)

type ctxKey struct{}

type Auth struct {
	db           *sql.DB
	cookieSecure bool
}

type Identity struct {
	ID   int
	Role string
}

func NewAuth(db *sql.DB, cookieSecure bool) *Auth {
	return &Auth{db: db, cookieSecure: cookieSecure}
}

func (a *Auth) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, fromCookie := extractToken(r)

		if token == "" {
			httpx.WriteError(w, http.StatusUnauthorized, errors.New("not logged in"))
			return
		}

		userID, role, err := session.Validate(a.db, token)

		if errors.Is(err, errs.ErrUnauthorized) {
			httpx.WriteError(w, http.StatusUnauthorized, errors.New("session expired"))
			return
		}

		if err != nil {
			log.Printf("auth validate: %v", err)
			httpx.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
			return
		}

		if err := session.Renew(a.db, token); err != nil {
			log.Printf("auth renew: %v", err)
		} else if fromCookie {
			session.SetCookie(w, token, session.SessionTTL, a.cookieSecure)
		}

		ctx := context.WithValue(r.Context(), ctxKey{}, Identity{ID: userID, Role: role})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// IdentityFromContext достаёт текущего пользователя, положенного Require.
func IdentityFromContext(ctx context.Context) (Identity, bool) {
	me, ok := ctx.Value(ctxKey{}).(Identity)
	return me, ok
}

func extractToken(r *http.Request) (string, bool) {
	if c, err := r.Cookie(session.CookieName); err == nil && c.Value != "" {
		return c.Value, true
	}

	if token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer "); ok && token != "" {
		return token, false
	}

	return "", false
}
