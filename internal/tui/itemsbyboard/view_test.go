package itemsbyboard_test

import (
	"strings"
	"testing"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/testutil"
	"github.com/rhajizada/donezo/internal/tui/boards"
	"github.com/rhajizada/donezo/internal/tui/itemsbyboard"
)

func TestItemsByBoardViewForListAndInputStates(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	ctx := testutil.MustContext()
	boardMenu := boards.New(ctx, svc)
	board, err := svc.CreateBoard(ctx, "Inbox")
	if err != nil {
		t.Fatalf("CreateBoard: %v", err)
	}
	boardMenu.List.SetItems(boards.NewList(&[]service.Board{*board}))
	boardMenu.List.Select(0)

	menu := itemsbyboard.New(ctx, svc, &boardMenu)
	menu.List.SetSize(80, 20)
	menu.List.Title = "Inbox"
	listView := menu.View().Content
	if strings.TrimSpace(listView) == "" {
		t.Fatalf("expected non-empty items-by-board list view")
	}

	menu.Context.State = itemsbyboard.CreateItemNameState
	menu.Input.SetValue("new item")
	inputView := menu.View().Content
	if !strings.Contains(inputView, "new item") {
		t.Fatalf("expected input text in input view")
	}
}
