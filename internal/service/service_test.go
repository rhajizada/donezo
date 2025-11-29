package service_test

import (
	"context"
	"testing"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/testutil"
)

func TestBoardLifecycle(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	ctx := testutil.MustContext()

	board := mustCreateBoard(ctx, t, svc, "Inbox")
	requirePositiveID(t, board.ID)

	boards := mustListBoards(ctx, t, svc)
	assertBoardNames(t, boards, []string{"Inbox"})

	board.Name = "Updated"
	updated := mustUpdateBoard(ctx, t, svc, board)
	requireEqualString(t, updated.Name, "Updated", "updated board name")

	mustDeleteBoard(ctx, t, svc, updated)

	afterDelete := mustListBoards(ctx, t, svc)
	assertBoardNames(t, afterDelete, []string{})
}

func TestItemLifecycle(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	ctx := testutil.MustContext()
	board := mustCreateBoard(ctx, t, svc, "Projects")

	item := mustCreateItem(ctx, t, svc, board, "task", "desc")
	itemsByBoard := mustListItemsByBoard(ctx, t, svc, board)
	requireEqualInt(t, len(*itemsByBoard), 1, "items by board")

	item.Tags = []string{"work", "go"}
	item.Completed = true
	item.Description = "updated desc"
	updated := mustUpdateItem(ctx, t, svc, item)
	requireEqualInt(t, len(updated.Tags), 2, "updated tags")
	requireTrue(t, updated.Completed, "item completed flag")
}

func TestTagCounts(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	ctx := testutil.MustContext()
	board := mustCreateBoard(ctx, t, svc, "Tags")
	item := mustCreateItem(ctx, t, svc, board, "tagged", "")

	item.Tags = []string{"work", "go"}
	item = mustUpdateItem(ctx, t, svc, item)

	requireEqualInt(t, mustCountItemsByTag(ctx, t, svc, "work"), 1, "count work tag")
	requireEqualInt(t, mustCountItemsByTag(ctx, t, svc, "go"), 1, "count go tag")

	item.Tags = []string{"go"}
	item = mustUpdateItem(ctx, t, svc, item)

	requireEqualInt(t, mustCountItemsByTag(ctx, t, svc, "work"), 0, "count work after removal")
	requireEqualInt(t, mustCountItemsByTag(ctx, t, svc, "go"), 1, "count go remains")

	mustDeleteItem(ctx, t, svc, item)
}

func mustCreateBoard(ctx context.Context, t *testing.T, svc *service.Service, name string) *service.Board {
	t.Helper()
	board, err := svc.CreateBoard(ctx, name)
	requireNoError(t, err, "CreateBoard")
	return board
}

func mustListBoards(ctx context.Context, t *testing.T, svc *service.Service) *[]service.Board {
	t.Helper()
	boards, err := svc.ListBoards(ctx)
	requireNoError(t, err, "ListBoards")
	return boards
}

func mustUpdateBoard(ctx context.Context, t *testing.T, svc *service.Service, board *service.Board) *service.Board {
	t.Helper()
	updated, err := svc.UpdateBoard(ctx, board)
	requireNoError(t, err, "UpdateBoard")
	return updated
}

func mustDeleteBoard(ctx context.Context, t *testing.T, svc *service.Service, board *service.Board) {
	t.Helper()
	requireNoError(t, svc.DeleteBoard(ctx, board), "DeleteBoard")
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
	requireNoError(t, err, "CreateItem")
	return item
}

func mustUpdateItem(ctx context.Context, t *testing.T, svc *service.Service, item *service.Item) *service.Item {
	t.Helper()
	updated, err := svc.UpdateItem(ctx, item)
	requireNoError(t, err, "UpdateItem")
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
	requireNoError(t, err, "ListItemsByBoard")
	return items
}

func mustCountItemsByTag(ctx context.Context, t *testing.T, svc *service.Service, tag string) int {
	t.Helper()
	count, err := svc.CountItemsByTag(ctx, tag)
	requireNoError(t, err, "CountItemsByTag")
	return int(count)
}

func mustDeleteItem(ctx context.Context, t *testing.T, svc *service.Service, item *service.Item) {
	t.Helper()
	requireNoError(t, svc.DeleteItem(ctx, item), "DeleteItem")
}

func requireNoError(t *testing.T, err error, context string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", context, err)
	}
}

func requirePositiveID(t *testing.T, id int64) {
	t.Helper()
	if id == 0 {
		t.Fatalf("expected ID to be set")
	}
}

func requireEqualString(t *testing.T, got, expected string, label string) {
	t.Helper()
	if got != expected {
		t.Fatalf("expected %s %q, got %q", label, expected, got)
	}
}

func requireEqualInt(t *testing.T, got, expected int, label string) {
	t.Helper()
	if got != expected {
		t.Fatalf("expected %s %d, got %d", label, expected, got)
	}
}

func requireTrue(t *testing.T, v bool, label string) {
	t.Helper()
	if !v {
		t.Fatalf("expected %s to be true", label)
	}
}

func assertBoardNames(t *testing.T, boards *[]service.Board, names []string) {
	t.Helper()
	if len(*boards) != len(names) {
		t.Fatalf("expected %d boards, got %d", len(names), len(*boards))
	}
	for i, b := range *boards {
		if b.Name != names[i] {
			t.Fatalf("unexpected board name %q at index %d", b.Name, i)
		}
	}
}
