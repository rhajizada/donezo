package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/testutil"
)

func TestBoardLifecycle(t *testing.T) {
	tests := []struct {
		name      string
		boardName string
		updated   string
	}{
		{name: "create update and delete board", boardName: "Inbox", updated: "Updated"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, cleanup := testutil.NewTestService(t)
			defer cleanup()

			ctx := testutil.MustContext()
			board := mustCreateBoard(ctx, t, svc, tt.boardName)
			require.Positive(t, board.ID)

			boards := mustListBoards(ctx, t, svc)
			assertBoardNames(t, boards, []string{tt.boardName})

			board.Name = tt.updated
			updated := mustUpdateBoard(ctx, t, svc, board)
			assert.Equal(t, tt.updated, updated.Name)

			mustDeleteBoard(ctx, t, svc, updated)

			afterDelete := mustListBoards(ctx, t, svc)
			assertBoardNames(t, afterDelete, []string{})
		})
	}
}

func TestItemLifecycle(t *testing.T) {
	tests := []struct {
		name            string
		boardName       string
		title           string
		description     string
		updatedDesc     string
		tags            []string
		expectCompleted bool
	}{
		{
			name:            "create and update item",
			boardName:       "Projects",
			title:           "task",
			description:     "desc",
			updatedDesc:     "updated desc",
			tags:            []string{"work", "go"},
			expectCompleted: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, cleanup := testutil.NewTestService(t)
			defer cleanup()

			ctx := testutil.MustContext()
			board := mustCreateBoard(ctx, t, svc, tt.boardName)

			item := mustCreateItem(ctx, t, svc, board, tt.title, tt.description)
			itemsByBoard := mustListItemsByBoard(ctx, t, svc, board)
			require.Len(t, *itemsByBoard, 1)

			item.Tags = tt.tags
			item.Completed = tt.expectCompleted
			item.Description = tt.updatedDesc
			updated := mustUpdateItem(ctx, t, svc, item)
			assert.Len(t, updated.Tags, len(tt.tags))
			assert.Equal(t, tt.updatedDesc, updated.Description)
			assert.Equal(t, tt.expectCompleted, updated.Completed)
		})
	}
}

func TestTagCounts(t *testing.T) {
	tests := []struct {
		name        string
		initialTags []string
		updatedTags []string
		wantBefore  map[string]int
		wantAfter   map[string]int
	}{
		{
			name:        "tag counts track updates",
			initialTags: []string{"work", "go"},
			updatedTags: []string{"go"},
			wantBefore:  map[string]int{"work": 1, "go": 1},
			wantAfter:   map[string]int{"work": 0, "go": 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, cleanup := testutil.NewTestService(t)
			defer cleanup()

			ctx := testutil.MustContext()
			board := mustCreateBoard(ctx, t, svc, "Tags")
			item := mustCreateItem(ctx, t, svc, board, "tagged", "")

			item.Tags = tt.initialTags
			item = mustUpdateItem(ctx, t, svc, item)

			for tag, want := range tt.wantBefore {
				assert.Equal(t, want, mustCountItemsByTag(ctx, t, svc, tag))
			}

			item.Tags = tt.updatedTags
			item = mustUpdateItem(ctx, t, svc, item)

			for tag, want := range tt.wantAfter {
				assert.Equal(t, want, mustCountItemsByTag(ctx, t, svc, tag))
			}

			mustDeleteItem(ctx, t, svc, item)
		})
	}
}

func mustCreateBoard(ctx context.Context, t *testing.T, svc *service.Service, name string) *service.Board {
	t.Helper()
	board, err := svc.CreateBoard(ctx, name)
	require.NoError(t, err)
	return board
}

func mustListBoards(ctx context.Context, t *testing.T, svc *service.Service) *[]service.Board {
	t.Helper()
	boards, err := svc.ListBoards(ctx)
	require.NoError(t, err)
	return boards
}

func mustUpdateBoard(ctx context.Context, t *testing.T, svc *service.Service, board *service.Board) *service.Board {
	t.Helper()
	updated, err := svc.UpdateBoard(ctx, board)
	require.NoError(t, err)
	return updated
}

func mustDeleteBoard(ctx context.Context, t *testing.T, svc *service.Service, board *service.Board) {
	t.Helper()
	require.NoError(t, svc.DeleteBoard(ctx, board))
}

func mustCreateItem(
	ctx context.Context,
	t *testing.T,
	svc *service.Service,
	board *service.Board,
	title, desc string,
) *service.Item {
	t.Helper()
	item, err := svc.CreateItem(ctx, board, title, desc)
	require.NoError(t, err)
	return item
}

func mustUpdateItem(ctx context.Context, t *testing.T, svc *service.Service, item *service.Item) *service.Item {
	t.Helper()
	updated, err := svc.UpdateItem(ctx, item)
	require.NoError(t, err)
	return updated
}

func mustListItemsByBoard(
	ctx context.Context,
	t *testing.T,
	svc *service.Service,
	board *service.Board,
) *[]service.Item {
	t.Helper()
	items, err := svc.ListItemsByBoard(ctx, board)
	require.NoError(t, err)
	return items
}

func mustCountItemsByTag(ctx context.Context, t *testing.T, svc *service.Service, tag string) int {
	t.Helper()
	count, err := svc.CountItemsByTag(ctx, tag)
	require.NoError(t, err)
	return int(count)
}

func mustDeleteItem(ctx context.Context, t *testing.T, svc *service.Service, item *service.Item) {
	t.Helper()
	require.NoError(t, svc.DeleteItem(ctx, item))
}

func assertBoardNames(t *testing.T, boards *[]service.Board, names []string) {
	t.Helper()
	require.Len(t, *boards, len(names))
	for i, b := range *boards {
		assert.Equal(t, names[i], b.Name)
	}
}
