package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/GuilhermeRossiKirsten/CustomerRegistryAPI/internal/config"
	"github.com/GuilhermeRossiKirsten/CustomerRegistryAPI/internal/database"
)

func main() {

	config.InitializeConfig()

	db, err := database.GetDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	row := db.QueryRow("SELECT 1+1")

	var result int

	err = row.Scan(&result)
	if err != nil {
		log.Fatalf("Failed to scan result: %v", err)
		os.Exit(1)
	}
	
	fmt.Printf("%v\n", result)
}
