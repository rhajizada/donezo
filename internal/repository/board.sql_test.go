package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rhajizada/donezo/internal/repository"
)

func TestBoardQueries(t *testing.T) {
	tests := []struct {
		name string
		run  func(*testing.T, *repository.Queries)
	}{
		{
			name: "create get list update and delete board",
			run: func(t *testing.T, q *repository.Queries) {
				ctx := context.Background()
				board, err := q.CreateBoard(ctx, "Inbox")
				require.NoError(t, err)
				assert.Positive(t, board.ID)
				assert.Equal(t, "Inbox", board.Name)

				fetched, err := q.GetBoardByID(ctx, board.ID)
				require.NoError(t, err)
				assert.Equal(t, board.ID, fetched.ID)
				assert.Equal(t, board.Name, fetched.Name)

				boards, err := q.ListBoards(ctx)
				require.NoError(t, err)
				require.Len(t, boards, 1)
				assert.Equal(t, board.ID, boards[0].ID)

				updated, err := q.UpdateBoardByID(ctx, repository.UpdateBoardByIDParams{
					ID:   board.ID,
					Name: "Projects",
				})
				require.NoError(t, err)
				assert.Equal(t, "Projects", updated.Name)

				require.NoError(t, q.DeleteBoardByID(ctx, board.ID))
				boards, err = q.ListBoards(ctx)
				require.NoError(t, err)
				assert.Empty(t, boards)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, q := newTestQueries(t)
			tt.run(t, q)
		})
	}
}
