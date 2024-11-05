package main

import (
	"flag"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rhajizada/donezo/internal/ui"
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

	p := tea.NewProgram(ui.NewModel(c), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Panic(err)
	}
}
