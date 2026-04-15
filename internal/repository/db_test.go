package repository_test

import (
	"context"
	"database/sql"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"

	"github.com/rhajizada/donezo/internal/repository"
)

func newTestQueries(t *testing.T) (*sql.DB, *repository.Queries) {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "repository-test.db")
	db, err := sql.Open("sqlite3", dbPath)
	require.NoError(t, err)

	baseFS, migrationsDir := migrationsFS(t)
	require.NoError(t, goose.SetDialect("sqlite3"))
	goose.SetBaseFS(baseFS)
	require.NoError(t, goose.Up(db, migrationsDir))

	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	return db, repository.New(db)
}

func mustBeginTx(t *testing.T, db *sql.DB) *sql.Tx {
	t.Helper()

	tx, err := db.BeginTx(context.Background(), nil)
	require.NoError(t, err)
	return tx
}

func mustCreateBoard(t *testing.T, q *repository.Queries, name string) repository.Board {
	t.Helper()

	board, err := q.CreateBoard(context.Background(), name)
	require.NoError(t, err)
	return board
}

func mustCreateItem(
	t *testing.T,
	q *repository.Queries,
	boardID int64,
	title string,
	description string,
) repository.Item {
	t.Helper()

	item, err := q.CreateItem(context.Background(), repository.CreateItemParams{
		BoardID:     boardID,
		Title:       title,
		Description: description,
	})
	require.NoError(t, err)
	return item
}

func tagsJSON(t *testing.T, value any) string {
	t.Helper()

	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		require.Failf(t, "unexpected tags type", "got %T", value)
		return ""
	}
}

func migrationsFS(t *testing.T) (fs.FS, string) {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok)

	root := filepath.Clean(filepath.Join(filepath.Dir(filename), "..", ".."))
	return os.DirFS(root), "data/sql/migrations"
}
