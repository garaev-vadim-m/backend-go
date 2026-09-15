package database

import (
	"log"
)

func RunMigrations() error {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		female VARCHAR(10),
		status_id INTEGER
	);
	`

	if _, err := DB.Exec(createUsersTable); err != nil {
		return err
	}

	log.Println("✓ Migration: users table created")
	return nil
}
