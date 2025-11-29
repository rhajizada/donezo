package boards

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/testutil"
)

func TestCopyBoardWritesMarkdown(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	ctx := testutil.MustContext()
	board, err := svc.CreateBoard(ctx, "Inbox")
	if err != nil {
		t.Fatalf("CreateBoard: %v", err)
	}

	_, err = svc.CreateItem(ctx, board, "task", "desc")
	if err != nil {
		t.Fatalf("CreateItem task: %v", err)
	}
	item2, err := svc.CreateItem(ctx, board, "done", "complete me")
	if err != nil {
		t.Fatalf("CreateItem done: %v", err)
	}
	item2.Completed = true
	if _, err := svc.UpdateItem(ctx, item2); err != nil {
		t.Fatalf("UpdateItem: %v", err)
	}

	items, err := svc.ListItemsByBoard(ctx, board)
	if err != nil {
		t.Fatalf("ListItemsByBoard: %v", err)
	}

	menu := New(ctx, svc)
	menu.List.SetItems(NewList(&[]service.Board{*board}))
	menu.List.Select(0)

	expected := service.ItemsToMarkdown(board.Name, *items)

	var captured []byte
	prevWrite := writeClipboardText
	writeClipboardText = func(data []byte) { captured = append([]byte{}, data...) }
	defer func() { writeClipboardText = prevWrite }()

	cmd := menu.Copy()
	if cmd != nil {
		cmd()
	}

	if string(captured) != expected {
		t.Fatalf("clipboard content mismatch\nwant:\n%s\n\ngot:\n%s", expected, string(captured))
	}
}

func TestBoardCreateRenameDeleteFlow(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	ctx := testutil.MustContext()
	menu := New(ctx, svc)

	// Create
	menu.InitCreateBoard()
	menu.Input.SetValue("Projects")
	model, cmd := menu.Update(tea.KeyMsg{Type: tea.KeyEnter})
	menu = model.(MenuModel)
	if cmd != nil {
		if msg := cmd(); msg != nil {
			model, _ = menu.Update(msg)
			menu = model.(MenuModel)
		}
	}
	if len(menu.List.Items()) != 1 {
		t.Fatalf("expected 1 board after create, got %d", len(menu.List.Items()))
	}
	created := menu.List.SelectedItem().(Item)
	if created.Board.Name != "Projects" {
		t.Fatalf("unexpected board name %q", created.Board.Name)
	}

	// Rename
	menu.InitRenameBoard()
	menu.Input.SetValue("Renamed")
	model, cmd = menu.Update(tea.KeyMsg{Type: tea.KeyEnter})
	menu = model.(MenuModel)
	if cmd != nil {
		if msg := cmd(); msg != nil {
			model, _ = menu.Update(msg)
			menu = model.(MenuModel)
		}
	}
	renamed := menu.List.SelectedItem().(Item)
	if renamed.Board.Name != "Renamed" {
		t.Fatalf("expected renamed board, got %q", renamed.Board.Name)
	}

	// Delete
	cmd = menu.DeleteBoard()
	msg := cmd()
	model, _ = menu.Update(msg)
	menu = model.(MenuModel)

	refresh := menu.ListBoards()
	msg = refresh()
	model, _ = menu.Update(msg)
	menu = model.(MenuModel)
	if len(menu.List.Items()) != 0 {
		t.Fatalf("expected board list empty after delete, got %d", len(menu.List.Items()))
	}
	boards, err := svc.ListBoards(ctx)
	if err != nil {
		t.Fatalf("ListBoards: %v", err)
	}
	if len(*boards) != 0 {
		t.Fatalf("expected no boards in service after delete, got %d", len(*boards))
	}
}
