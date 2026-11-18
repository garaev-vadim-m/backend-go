package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"hello-go-backend/internal/auth"
)

// SeedAdmin создает роли admin и user и технического админа, если их еще нет.
// Пароль админа хешируется, как и у остальных пользователей.
func SeedAdmin() error {
	email := getenvDefault("ADMIN_EMAIL", "email@email")
	password := getenvDefault("ADMIN_PASSWORD", "1235")
	name := getenvDefault("ADMIN_NAME", "tech")

	adminRuleID, err := ensureRule("admin", "Администратор")
	if err != nil {
		return err
	}
	if _, err := ensureRule("user", "Пользователь"); err != nil {
		return err
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return fmt.Errorf("seed admin hash: %w", err)
	}

	var stored sql.NullString
	err = DB.QueryRow(`SELECT password FROM users WHERE email = $1`, email).Scan(&stored)

	if errors.Is(err, sql.ErrNoRows) {
		_, err = DB.Exec(
			`INSERT INTO users (name, email, password, rule_id) VALUES ($1, $2, $3, $4)`,
			name, email, hash, adminRuleID,
		)
		if err != nil {
			return fmt.Errorf("seed admin insert: %w", err)
		}
		log.Printf("✓ admin user created: %s", email)
		return nil
	}
	if err != nil {
		return fmt.Errorf("seed admin check: %w", err)
	}

	if !stored.Valid || !auth.CheckPassword(stored.String, password) {
		_, err = DB.Exec(
			`UPDATE users SET password = $1, rule_id = $2 WHERE email = $3`,
			hash, adminRuleID, email,
		)
		if err != nil {
			return fmt.Errorf("seed admin update: %w", err)
		}
		log.Printf("✓ admin password updated: %s", email)
		return nil
	}

	_, err = DB.Exec(
		`UPDATE users SET rule_id = $1 WHERE email = $2 AND rule_id IS DISTINCT FROM $1`,
		adminRuleID, email,
	)
	if err != nil {
		return fmt.Errorf("seed admin rule update: %w", err)
	}

	return nil
}

// ensureRule находит роль по code или создает её.
func ensureRule(code, name string) (int, error) {
	var id int
	err := DB.QueryRow(`SELECT id FROM rules WHERE code = $1 ORDER BY id LIMIT 1`, code).Scan(&id)
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("seed rule %s check: %w", code, err)
	}

	err = DB.QueryRow(`INSERT INTO rules (name, code) VALUES ($1, $2) RETURNING id`, name, code).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("seed rule %s insert: %w", code, err)
	}
	return id, nil
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
