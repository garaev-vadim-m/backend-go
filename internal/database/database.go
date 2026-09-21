package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error

	var db *sql.DB

	for i := range 5 {
		db, err = sql.Open("postgres", connStr)

		if err != nil {
			return err
		}

		// Проверяем подключение
		if err := db.Ping(); err == nil {
			DB = db
			fmt.Println("✓ Connected to PostgreSQL")
			return nil
		}

		fmt.Printf("Attempt %d failed, retrying in 2 seconds...\n", i+1)
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
}

func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
