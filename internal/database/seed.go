package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

// SeedAdmin создает технического админа, если его еще нет.
// Логин и пароль берутся из ADMIN_EMAIL / ADMIN_PASSWORD, по умолчанию email@email / 1235.
func SeedAdmin() error {
	email := getenvDefault("ADMIN_EMAIL", "email@email")
	password := getenvDefault("ADMIN_PASSWORD", "1235")
	name := getenvDefault("ADMIN_NAME", "tech")

	ruleID, err := ensureAdminRule()
	if err != nil {
		return err
	}

	var exists bool
	err = DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, email).Scan(&exists)
	if err != nil {
		return fmt.Errorf("seed admin check: %w", err)
	}

	if !exists {
		_, err = DB.Exec(
			`INSERT INTO users (name, email, password, rule_id) VALUES ($1, $2, $3, $4)`,
			name, email, password, ruleID,
		)
		if err != nil {
			return fmt.Errorf("seed admin insert: %w", err)
		}
		log.Printf("✓ admin user created: %s", email)
		return nil
	}

	result, err := DB.Exec(
		`UPDATE users
		 SET rule_id = $1, password = $2
		 WHERE email = $3
		   AND (rule_id IS DISTINCT FROM $1 OR password IS DISTINCT FROM $2)`,
		ruleID, password, email,
	)
	if err != nil {
		return fmt.Errorf("seed admin update: %w", err)
	}
	if affected, err := result.RowsAffected(); err == nil && affected > 0 {
		log.Printf("✓ admin user updated: %s", email)
	}

	return nil
}

func ensureAdminRule() (int, error) {
	var id int
	err := DB.QueryRow(`
		SELECT id FROM rules
		WHERE lower(code) = 'admin' OR lower(name) = 'admin'
		ORDER BY id LIMIT 1`,
	).Scan(&id)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		return 0, fmt.Errorf("seed admin rule check: %w", err)
	}

	err = DB.QueryRow(`INSERT INTO rules (name, code) VALUES ('admin', 'admin') RETURNING id`).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("seed admin rule insert: %w", err)
	}
	return id, nil
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
