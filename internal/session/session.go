package session

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"hello-go-backend/internal/errs"
	"time"
)

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func Create(db *sql.DB, userID int, ttl time.Duration) (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	token := hex.EncodeToString(b)

	_, err := db.Exec(
		`INSERT INTO session (user_id, token_hash, expire_at) VALUES ($1 , $2, $3)`, userID, hashToken(token), time.Now().Add(ttl),
	)

	if err != nil {
		return "", nil
	}

	return token, nil
}

// Validate возвращает id пользователя, если сессия активна.
// Этим будет пользоваться middleware.
func Validate(db *sql.DB, token string) (int, error) {
	var userID int
	err := db.QueryRow(`
		SELECT user_id FROM sessions
		WHERE token_hash = $1
		  AND revoked_at IS NULL
		  AND expires_at > now()`,
		hashToken(token),
	).Scan(&userID)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, errs.ErrUnauthorized
	}
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// Revoke помечает сессию завершённой. Строка остаётся для истории.
// Повторный выход ошибки не даёт.
func Revoke(db *sql.DB, token string) error {
	_, err := db.Exec(`
		UPDATE sessions SET revoked_at = now()
		WHERE token_hash = $1 AND revoked_at IS NULL`,
		hashToken(token),
	)
	return err
}
