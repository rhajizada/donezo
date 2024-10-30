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
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/mattn/go-sqlite3"
)

// @title donezo API
// @description Swagger API documentation for donezo.
// @version 0.1.0

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	configPath := flag.String("config", "/etc/donezo/config.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadServerConfig(*configPath)
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
	rq := repository.New(db)

	// Create handler
	h := handler.New(rq, cfg.JWT.Secret, cfg.JWT.Expiration)

	// Register API routes
	r := http.NewServeMux()
	authMiddleware := middleware.AuthMiddleware(cfg.JWT.Secret)
	boards := authMiddleware(router.RegisterBoardRoutes(h))
	token := router.RegisterTokenRoutes(h)
	r.Handle("/api/boards/", http.StripPrefix("/api", boards))
	r.Handle("/api/token/", http.StripPrefix("/api", token))
	r.Handle("/swagger/", httpSwagger.WrapHandler)
	r.HandleFunc("/healthz/", handler.Healthz)
	lm := middleware.Logging(r)

	// Start the server
	log.Printf("Server is running on port %v\n", cfg.Port)
	addr := fmt.Sprintf(":%v", cfg.Port)
	if err := http.ListenAndServe(addr, lm); err != nil {
		log.Panicf("Could not start server: %s\n", err.Error())
	}
}
