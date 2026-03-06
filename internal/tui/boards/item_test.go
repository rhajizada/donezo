package boards_test

import (
	"testing"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/testutil"
	"github.com/rhajizada/donezo/internal/tui/boards"
)

func TestBoardItemAccessors(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	ctx := testutil.MustContext()
	board, err := svc.CreateBoard(ctx, "Inbox")
	if err != nil {
		t.Fatalf("CreateBoard: %v", err)
	}

	item := boards.NewItem(board).(boards.Item)
	if item.Title() != "Inbox" {
		t.Fatalf("expected title Inbox, got %q", item.Title())
	}
	if item.FilterValue() != "Inbox" {
		t.Fatalf("expected filter value Inbox, got %q", item.FilterValue())
	}
	if item.Description() == "" {
		t.Fatalf("expected non-empty description")
	}

	list := boards.NewList(&[]service.Board{*board})
	if len(list) != 1 {
		t.Fatalf("expected one list item, got %d", len(list))
	}
}
