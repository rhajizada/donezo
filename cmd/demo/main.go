package main

import (
	"flag"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/rhajizada/donezo/pkg/client"
	"github.com/rhajizada/donezo/pkg/config"
)

func main() {
	defaultConfigPath, err := config.GetDefaultConfigPath()
	if err != nil {
		log.Panic(err)
	}
	// Define the config flag with the determined default path
	configPath := flag.String("config", defaultConfigPath, "Path to configuration file")
	flag.Parse()
	// Load configuration
	cfg, err := config.LoadClientConfig(*configPath)
	if err != nil {
		log.Panicf("Error loading config: %v", err)
	}

	// Initialize the faker with a seed for reproducibility (optional)

	c := client.New(
		cfg.BaseURL,
		cfg.ApiToken,
		time.Second*5,
	)

	if err := c.Healthy(); err != nil {
		log.Panicf("Cannot connect to %s: %v", c.BaseURL, err)
	}

	err = c.ValidateToken()
	if err != nil {
		log.Panic(err)
	}

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
			_, err := c.CreateItem(board, title, description)
			if err != nil {
				log.Fatalf("Failed to create an item: %v", err)
			}
		}
	}
}
