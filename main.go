package main

import (
	"context"
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pressly/goose/v3"
	"github.com/rhajizada/donezo/internal/repository"
	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/tui/app"
	"golang.design/x/clipboard"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed data/sql/migrations/*.sql
var migrations embed.FS
var Version = "dev"

func main() {
	versionFlag := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	// If the version flag is provided, print version info and exit.
	if *versionFlag {
		fmt.Printf("donezo %s\n", Version)
		os.Exit(0)
	}
	err := clipboard.Init()
	if err != nil {
		log.Panicf("unable to access system clipboard: %v", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Panicf("unable to determine user home directory: %v", err)
	}
	donezoDir := filepath.Join(homeDir, ".donezo")
	if _, err = os.Stat(donezoDir); os.IsNotExist(err) {
		err = os.Mkdir(donezoDir, 0700)
		if err != nil {
			log.Panicf("failed to create directory %s: %v", donezoDir, err)
		}
	}

	dbPath := filepath.Join(donezoDir, "data.db")

	// Load database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Panicf("failed to open database %s: %v", dbPath, err)
	}
	defer db.Close()

	// Set Goose dialect to SQLite
	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Panicf("failed to set Goose dialect: %v", err)
	}

	// Set the embedded migrations as the base FS for Goose
	goose.SetBaseFS(migrations)

	// Apply all up migrations
	if err := goose.Up(db, "data/sql/migrations"); err != nil {
		log.Panicf("failed to apply migrations: %v", err)
	}

	r := repository.New(db)
	s := service.New(r)
	ctx := context.Background()

	m := app.New(ctx, s)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Panicf("error running program: %v", err)
	}
}
