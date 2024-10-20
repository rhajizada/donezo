package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/pressly/goose/v3"
	"github.com/rhajizada/donezo/internal/repository"
)

func main() {
	// Define paths
	projectRoot, err := getProjectRoot()
	if err != nil {
		log.Fatalf("Failed to determine project root: %v", err)
	}

	dbPath := filepath.Join(projectRoot, "db.sqlite")
	migrationsDir := filepath.Join(projectRoot, "data", "sql", "migrations")

	// Open the SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Could not open database: %v", err)
	}
	defer db.Close()

	// Verify the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	// Initialize Goose
	goose.SetDialect("sqlite3")

	// Run migrations
	if err := goose.Up(db, migrationsDir); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	ctx := context.Background()

	fmt.Println("Migrations applied successfully.")
	repo := repository.New(db)
	board, err := repo.CreateBoard(ctx, "Agenda")
	if err != nil {
		msg := fmt.Sprintf("failed to create table 'Agenda': %v", err)
		panic(msg)
	} else {
		fmt.Printf("Created board 'Agenda' with id %d", board.ID)
	}
	item, err := repo.CreateItem(ctx, repository.CreateItemParams{
		board.ID,
		"Donezo",
		"Cool todo app",
	})

	if err != nil {
		msg := fmt.Sprintf("failed to add item to 'Agenda': %v", err)
		panic(msg)
	} else {
		fmt.Printf("Added item %d to 'Agenda'\n", item.ID)
	}

	repo.UpdateItemByID(ctx, repository.UpdateItemByIDParams{
		item.Title,
		item.Description,
		1,
		item.ID,
	})
}

// getProjectRoot returns the absolute path to the project root.
func getProjectRoot() (string, error) {
	// Assuming main.go is in cmd/, so project root is one level up
	cwd, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		return "", err
	}

	// Navigate up to the project root
	projectRoot := cwd
	// Alternatively, you can use a more reliable method to find the project root
	// For simplicity, we're assuming the current working directory is the project root
	return projectRoot, nil
}
