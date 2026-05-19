package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitializeConfig() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		return
	}
}

func Env(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	log.Printf("Environment variable %s not set, using default value: %s", key, defaultVal)
	return defaultVal
}
