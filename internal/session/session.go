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

const (
	SessionTTL  = 24 * time.Hour
	RenewBefore = 12 * time.Hour
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
		`INSERT INTO sessions (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`, userID, hashToken(token), time.Now().Add(ttl),
	)

	if err != nil {
		return "", err
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

// Renew продлевает сессию, если до её конца осталось меньше RenewBefore.
// Вызывается из middleware после Validate.
func Renew(db *sql.DB, token string) error {
	_, err := db.Exec(`
		UPDATE sessions
		SET expires_at = now() + make_interval(secs => $2)
		WHERE token_hash = $1
		  AND revoked_at IS NULL
		  AND expires_at > now()
		  AND expires_at < now() + make_interval(secs => $3)`,
		hashToken(token), SessionTTL.Seconds(), RenewBefore.Seconds(),
	)
	return err
}

// Cleanup удаляет истекшие сессии без выхода, старше keepFor.
// Отозванные (revoked_at) не трогаем: это история входов и выходов.
func Cleanup(db *sql.DB, keepFor time.Duration) (int64, error) {
	res, err := db.Exec(`
		DELETE FROM sessions
		WHERE revoked_at IS NULL
		  AND expires_at < now() - make_interval(secs => $1)`,
		keepFor.Seconds(),
	)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
