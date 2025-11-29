package testutil

import (
	"context"
	"database/sql"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	_ "github.com/mattn/go-sqlite3" // sqlite driver
	"github.com/pressly/goose/v3"
	"github.com/rhajizada/donezo/internal/repository"
	"github.com/rhajizada/donezo/internal/service"
)

// MustContext returns a background context for tests.
func MustContext() context.Context {
	return context.Background()
}

// NewTestService spins up a temporary SQLite database, runs migrations,
// and returns a ready-to-use Service plus a cleanup function.
func NewTestService(t *testing.T) (*service.Service, func()) {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}

	baseFS, migrationsDir := migrationsFS(t)

	err = goose.SetDialect("sqlite3")
	if err != nil {
		t.Fatalf("SetDialect: %v", err)
	}
	goose.SetBaseFS(baseFS)

	err = goose.Up(db, migrationsDir)
	if err != nil {
		t.Fatalf("goose.Up: %v", err)
	}

	repo := repository.New(db)
	svc := service.New(repo)

	return svc, func() {
		_ = db.Close()
	}
}

func migrationsFS(t *testing.T) (fs.FS, string) {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("runtime.Caller failed")
	}

	root := filepath.Clean(filepath.Join(filepath.Dir(filename), "..", ".."))
	return os.DirFS(root), "data/sql/migrations"
}
