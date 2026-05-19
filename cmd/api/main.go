// @title Customer Registry API
// @version 1.0
// @description Customer registration and management API
// @host localhost:8080
// @basePath /
// @schemes http https
package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/GuilhermeRossiKirsten/CustomerRegistryAPI/docs"
	"github.com/GuilhermeRossiKirsten/CustomerRegistryAPI/internal/config"
	"github.com/GuilhermeRossiKirsten/CustomerRegistryAPI/internal/customer"
	"github.com/GuilhermeRossiKirsten/CustomerRegistryAPI/internal/database"
	"github.com/GuilhermeRossiKirsten/CustomerRegistryAPI/internal/health"
)

func main() {

	config.InitializeConfig()

	db, err := database.GetDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	repo := customer.NewRepository(db)
	service := customer.NewService(repo)
	handler := customer.NewHandler(service)

	mux := http.NewServeMux()
	handler.Register(mux)
	health.NewHandler(db).Register(mux)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	addr := ":" + config.Env("APP_PORT", "8080")
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
