package database

import (
	"log"
)

func RunMigrations() error {
	createStatusTable := `
	CREATE TABLE IF NOT EXISTS status (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL
	);
	`

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		female VARCHAR(10),
		status_id INTEGER REFERENCES status(id)
	);
	`

	log.Println("Creating status table...")
	if _, err := DB.Exec(createStatusTable); err != nil {
		log.Printf("Error creating status table: %v", err)
		return err
	}
	log.Println("✓ Status table created")

	log.Println("Creating users table...")
	if _, err := DB.Exec(createUsersTable); err != nil {
		log.Printf("Error creating users table: %v", err)
		return err
	}
	log.Println("✓ Users table created")

	return nil
}
