package session

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
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

func Delete(db *sql.DB, token string) error {
	_, err := db.Exec(`DELETE FROM sessions WHERE token_hash = $1`, hashToken(token))
	return err
}
