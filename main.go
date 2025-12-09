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

var Version = "dev" //nolint:gochecknoglobals // overridden at build time via ldflags

func main() {
	if err := run(); err != nil {
		log.Panic(err)
	}
}

func run() error {
	versionFlag := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Fprintf(os.Stdout, "donezo %s\n", Version)
		return nil
	}

	if err := clipboard.Init(); err != nil {
		return fmt.Errorf("unable to access system clipboard: %w", err)
	}

	dbPath, err := ensureDataDir()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database %s: %w", dbPath, err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			log.Printf("failed to close database %s: %v", dbPath, cerr)
		}
	}()

	if migrateErr := runMigrations(db); migrateErr != nil {
		return migrateErr
	}

	r := repository.New(db)
	s := service.New(r)
	ctx := context.Background()

	m := app.New(ctx, s)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, programErr := p.Run(); programErr != nil {
		return fmt.Errorf("error running program: %w", programErr)
	}

	return nil
}

func ensureDataDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to determine user home directory: %w", err)
	}

	donezoDir := filepath.Join(homeDir, ".donezo")
	if _, err = os.Stat(donezoDir); err != nil {
		if os.IsNotExist(err) {
			if mkErr := os.Mkdir(donezoDir, 0700); mkErr != nil {
				return "", fmt.Errorf("failed to create directory %s: %w", donezoDir, mkErr)
			}
		} else {
			return "", fmt.Errorf("failed to check directory %s: %w", donezoDir, err)
		}
	}

	return filepath.Join(donezoDir, "data.db"), nil
}

func runMigrations(db *sql.DB) error {
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set Goose dialect: %w", err)
	}

	goose.SetBaseFS(migrations)

	if err := goose.Up(db, "data/sql/migrations"); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
