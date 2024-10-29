package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rhajizada/donezo/client"
	"github.com/rhajizada/donezo/internal/config"
	"github.com/rhajizada/donezo/internal/ui"
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
		log.Panicf("Error loading config: %v", err)
	}

	// Initialize the faker with a seed for reproducibility (optional)

	cli := client.New(
		cfg.BaseURL,
		cfg.ApiToken,
		cfg.Duration,
	)

	if err := cli.Healthy(); err != nil {
		log.Panicf("Cannot connect to %s: %v", cli.BaseURL, err)
	}

	p := tea.NewProgram(ui.NewModel(cli), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
