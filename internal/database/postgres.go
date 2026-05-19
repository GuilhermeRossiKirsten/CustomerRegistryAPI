package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/GuilhermeRossiKirsten/CustomerRegistryAPI/internal/config"
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
	host := config.Env("DB_HOST", "localhost")
	port := config.Env("DB_PORT", "5432")
	user := config.Env("DB_USER", "app")
	dbname := config.Env("DB_NAME", "customers-registry")
	password := config.Env("DB_PASSWORD", "app")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}
