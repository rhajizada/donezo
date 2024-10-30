package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/rhajizada/donezo/internal/config"
	"github.com/rhajizada/donezo/pkg/client"
)

func main() {
	var defaultConfigPath string

	// Check if XDG_CONFIG_HOME is set
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		defaultConfigPath = filepath.Join(xdgConfigHome, "donezo", "config.yaml")
	} else {
		// Fallback to ~/.config
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Unable to determine user home directory: %v", err)
		}
		defaultConfigPath = filepath.Join(homeDir, ".config", "donezo", "config.yaml")
	}

	// Define the config flag with the determined default path
	configPath := flag.String("config", defaultConfigPath, "Path to configuration file")
	flag.Parse()
	// Load configuration
	cfg, err := config.LoadClientConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize the faker with a seed for reproducibility (optional)
	gofakeit.Seed(time.Now().UnixNano())

	c := client.New(
		cfg.BaseURL,
		cfg.ApiToken,
		time.Second*15,
	)

	for i := 1; i < 10; i++ {
		// Generate a random sentence for the board name
		name := gofakeit.Sentence(3) // Generates a sentence with approximately 3 words
		board, err := c.CreateBoard(name)
		if err != nil {
			log.Fatalf("Failed to create a board: %v", err)
		}

		for j := 1; j < 10; j++ {
			// Generate random sentences for item title and description
			title := gofakeit.Sentence(2)       // Approximately 2 words
			description := gofakeit.Sentence(6) // Approximately 6 words
			_, err := c.AddItem(board, title, description)
			if err != nil {
				log.Fatalf("Failed to create an item: %v", err)
			}
		}
	}
}
