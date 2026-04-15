package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rhajizada/donezo/internal/repository"
)

func TestItemQueries(t *testing.T) {
	tests := []struct {
		name string
		run  func(*testing.T, *sql.DB, *repository.Queries)
	}{
		{
			name: "create get list update and delete item",
			run: func(t *testing.T, _ *sql.DB, q *repository.Queries) {
				ctx := context.Background()
				board := mustCreateBoard(t, q, "Inbox")
				item := mustCreateItem(t, q, board.ID, "task", "desc")

				require.NoError(t, q.AddTagToItemByID(ctx, repository.AddTagToItemByIDParams{
					ItemID: item.ID,
					Tag:    "work",
				}))
				require.NoError(t, q.AddTagToItemByID(ctx, repository.AddTagToItemByIDParams{
					ItemID: item.ID,
					Tag:    "go",
				}))

				fetched, err := q.GetItemByID(ctx, item.ID)
				require.NoError(t, err)
				assert.Equal(t, item.ID, fetched.ID)
				assert.Equal(t, "task", fetched.Title)
				assert.Contains(t, tagsJSON(t, fetched.Tags), "work")

				items, err := q.ListItemsByBoardID(ctx, board.ID)
				require.NoError(t, err)
				require.Len(t, items, 1)
				assert.Equal(t, item.ID, items[0].ID)

				updated, err := q.UpdateItemByID(ctx, repository.UpdateItemByIDParams{
					ID:          item.ID,
					Title:       "renamed",
					Description: "updated",
					Completed:   true,
				})
				require.NoError(t, err)
				assert.Equal(t, "renamed", updated.Title)
				assert.True(t, updated.Completed)

				require.NoError(t, q.DeleteItemByID(ctx, item.ID))
				items, err = q.ListItemsByBoardID(ctx, board.ID)
				require.NoError(t, err)
				assert.Empty(t, items)
			},
		},
		{
			name: "with tx uses transaction-scoped queries",
			run: func(t *testing.T, db *sql.DB, q *repository.Queries) {
				ctx := context.Background()
				board := mustCreateBoard(t, q, "Inbox")
				tx := mustBeginTx(t, db)

				txQueries := q.WithTx(tx)
				_, err := txQueries.CreateItem(ctx, repository.CreateItemParams{
					BoardID:     board.ID,
					Title:       "tx item",
					Description: "inside tx",
				})
				require.NoError(t, err)

				require.NoError(t, tx.Commit())

				items, err := q.ListItemsByBoardID(ctx, board.ID)
				require.NoError(t, err)
				require.Len(t, items, 1)
				assert.Equal(t, "tx item", items[0].Title)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, q := newTestQueries(t)
			tt.run(t, db, q)
		})
	}
}
