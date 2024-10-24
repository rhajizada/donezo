package main

import (
	"flag"
	"log"
	"time"

	"github.com/rhajizada/donezo/internal/auth"
	"github.com/rhajizada/donezo/internal/config"
)

func main() {
	configPath := flag.String("config", "/etc/donezo/config.yaml", "Path to configuration file")
	expiration := flag.Duration("expiration", 24*time.Hour, "Token duration")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	token, err := auth.GenerateJWT(cfg.JWTSecret, "donezo", *expiration)
	if err != nil {
		log.Fatalf("Failed to generate JWT token: %v", err)
	}
	log.Printf("Token: %s", token)
}
