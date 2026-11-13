package database

import (
	"log"
)

func RunMigrations() error {
	createStatusTable := `
	CREATE TABLE IF NOT EXISTS status (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		code VARCHAR(255) NOT NULL
	);
	`

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		surname VARCHAR(255),
		female VARCHAR(10),
		email VARCHAR(255),
		status_id INTEGER REFERENCES user_status(id),
		rule_id INTEGER REFERENCES rules(id)
	);
	`

	alterUsersTable := `
	ALTER TABLE users ADD COLUMN IF NOT EXISTS surname VARCHAR(255);
	ALTER TABLE users ADD COLUMN IF NOT EXISTS email VARCHAR(255);
	`

	createUserStatusTable := `
	CREATE TABLE IF NOT EXISTS user_status (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		code VARCHAR(255) NOT NULL
	);
	`

	createRulesTable := `
	CREATE TABLE IF NOT EXISTS rules (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		code VARCHAR(255) NOT NULL
	);
	`

	createProductsTable := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description VARCHAR(255),
		status_id INTEGER REFERENCES status(id)
	);
	`

	// Порядок важен! Сначала таблицы без зависимостей
	tables := []struct {
		name string
		sql  string
	}{
		{"status", createStatusTable},
		{"rules", createRulesTable},
		{"user_status", createUserStatusTable},
		{"users", createUsersTable},
		{"users_alert", alterUsersTable},
		{"products", createProductsTable},
	}

	for _, table := range tables {
		log.Printf("Creating %s table...", table.name)

		if _, err := DB.Exec(table.sql); err != nil {
			log.Printf("Error creating %s table: %v", table.name, err)
			return err
		}

		log.Printf("✓ %s table created", table.name)
	}

	return nil
}
