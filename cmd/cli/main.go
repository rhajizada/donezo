package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rhajizada/donezo/internal/config"
	"github.com/rhajizada/donezo/internal/ui"
	"github.com/rhajizada/donezo/pkg/client"
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

	cli := client.New(
		cfg.BaseURL,
		cfg.ApiToken,
		time.Second*5,
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
