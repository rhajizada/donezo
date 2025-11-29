package service_test

import (
	"testing"

	"github.com/rhajizada/donezo/internal/testutil"
)

func TestBoardLifecycle(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	ctx := testutil.MustContext()

	board, err := svc.CreateBoard(ctx, "Inbox")
	if err != nil {
		t.Fatalf("CreateBoard: %v", err)
	}
	if board.ID == 0 {
		t.Fatalf("expected board ID to be set")
	}

	boards, err := svc.ListBoards(ctx)
	if err != nil {
		t.Fatalf("ListBoards: %v", err)
	}
	if len(*boards) != 1 || (*boards)[0].Name != "Inbox" {
		t.Fatalf("unexpected boards: %+v", boards)
	}

	board.Name = "Updated"
	updated, err := svc.UpdateBoard(ctx, board)
	if err != nil {
		t.Fatalf("UpdateBoard: %v", err)
	}
	if updated.Name != "Updated" {
		t.Fatalf("expected updated board name, got %s", updated.Name)
	}

	if err := svc.DeleteBoard(ctx, updated); err != nil {
		t.Fatalf("DeleteBoard: %v", err)
	}

	afterDelete, err := svc.ListBoards(ctx)
	if err != nil {
		t.Fatalf("ListBoards after delete: %v", err)
	}
	if len(*afterDelete) != 0 {
		t.Fatalf("expected no boards after delete, got %d", len(*afterDelete))
	}
}

func TestItemAndTagLifecycle(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	ctx := testutil.MustContext()
	board, err := svc.CreateBoard(ctx, "Projects")
	if err != nil {
		t.Fatalf("CreateBoard: %v", err)
	}

	item, err := svc.CreateItem(ctx, board, "task", "desc")
	if err != nil {
		t.Fatalf("CreateItem: %v", err)
	}

	itemsByBoard, err := svc.ListItemsByBoard(ctx, board)
	if err != nil {
		t.Fatalf("ListItemsByBoard: %v", err)
	}
	if len(*itemsByBoard) != 1 {
		t.Fatalf("expected 1 item, got %d", len(*itemsByBoard))
	}

	// Add tags and mark complete.
	item.Tags = []string{"work", "go"}
	item.Completed = true
	item.Description = "updated desc"
	updated, err := svc.UpdateItem(ctx, item)
	if err != nil {
		t.Fatalf("UpdateItem: %v", err)
	}
	if len(updated.Tags) != 2 {
		t.Fatalf("expected tags to persist, got %+v", updated.Tags)
	}
	if !updated.Completed {
		t.Fatalf("expected item to be completed")
	}

	byTag, err := svc.ListItemsByTag(ctx, "work")
	if err != nil {
		t.Fatalf("ListItemsByTag: %v", err)
	}
	if len(*byTag) != 1 {
		t.Fatalf("expected 1 item for tag 'work', got %d", len(*byTag))
	}

	count, err := svc.CountItemsByTag(ctx, "work")
	if err != nil {
		t.Fatalf("CountItemsByTag: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 item counted for 'work', got %d", count)
	}

	// Remove a tag and ensure counts adjust.
	item.Tags = []string{"go"}
	_, err = svc.UpdateItem(ctx, item)
	if err != nil {
		t.Fatalf("UpdateItem removing tag: %v", err)
	}

	countWork, err := svc.CountItemsByTag(ctx, "work")
	if err != nil {
		t.Fatalf("CountItemsByTag work: %v", err)
	}
	if countWork != 0 {
		t.Fatalf("expected work tag to be removed, count=%d", countWork)
	}

	countGo, err := svc.CountItemsByTag(ctx, "go")
	if err != nil {
		t.Fatalf("CountItemsByTag go: %v", err)
	}
	if countGo != 1 {
		t.Fatalf("expected go tag to remain, count=%d", countGo)
	}

	if err := svc.DeleteItem(ctx, item); err != nil {
		t.Fatalf("DeleteItem: %v", err)
	}
}
