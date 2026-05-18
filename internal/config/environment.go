package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func InitializeConfig() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		return
	}
}
