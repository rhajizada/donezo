package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pressly/goose"
	"github.com/rhajizada/donezo/internal/config"
	"github.com/rhajizada/donezo/internal/handler"
	"github.com/rhajizada/donezo/internal/middleware"
	"github.com/rhajizada/donezo/internal/repository"
	"github.com/rhajizada/donezo/internal/router"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	configPath := flag.String("config", "/app/config.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Load database
	db, err := sql.Open("sqlite3", cfg.Database)
	if err != nil {
		log.Fatalf("Failed to open database %s: %v", cfg.Database, err)
	}
	defer db.Close()

	// Ensure the migrations directory exists
	migrationsDir := "data/sql/migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		log.Fatalf("Migrations directory does not exist: %s", migrationsDir)
	}

	// Set Goose dialect to SQLite
	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatalf("Failed to set Goose dialect: %v", err)
	}

	// Apply all up migrations
	if err := goose.Up(db, migrationsDir); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	// Create repository
	r := repository.New(db)

	// Create handler
	h := handler.New(r)

	// Register API routes
	router := router.RegisterApiRoutes(h)

	// Add middlewares
	mdlwr := middleware.CreateStack(middleware.Logging, middleware.AuthMiddleware(cfg.JWTSecret))

	// Start the server
	log.Printf("Server is running on port %v\n", cfg.Port)
	addr := fmt.Sprintf(":%v", cfg.Port)
	if err := http.ListenAndServe(addr, mdlwr(router)); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
