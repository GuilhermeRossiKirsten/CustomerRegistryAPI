package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func GetDBConnection() (*sql.DB, error) {
	db, err := sql.Open("pgx", getDBConnectionStr())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil, err
	}

	if err := verifyDBConnection(db); err != nil {
		log.Fatalf("Database connection verification failed: %v", err)
		return nil, err
	}

	return db, nil
}

func verifyDBConnection(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
		return err
	}
	return nil
}

func getDBConnectionStr() string {
	host := env("DB_HOST", "localhost")
	port := env("DB_PORT", "5432")
	user := env("DB_USER", "app")
	dbname := env("DB_NAME", "customers-registry")
	password := env("DB_PASSWORD", "app")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

func env(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	log.Printf("Environment variable %s not set, using default value: %s", key, defaultVal)
	return defaultVal
}
