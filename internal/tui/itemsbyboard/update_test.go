package itemsbyboard

import (
	"encoding/json"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/rhajizada/donezo/internal/repository"
	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/testutil"
	"github.com/rhajizada/donezo/internal/tui/boards"
)

func TestCopySavesItemJSON(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	ctx := testutil.MustContext()
	board, err := svc.CreateBoard(ctx, "Inbox")
	if err != nil {
		t.Fatalf("CreateBoard: %v", err)
	}
	item, err := svc.CreateItem(ctx, board, "title", "desc")
	if err != nil {
		t.Fatalf("CreateItem: %v", err)
	}
	item.Tags = []string{"tag1", "tag2"}
	item.Completed = true
	item, err = svc.UpdateItem(ctx, item)
	if err != nil {
		t.Fatalf("UpdateItem: %v", err)
	}

	items, err := svc.ListItemsByBoard(ctx, board)
	if err != nil {
		t.Fatalf("ListItemsByBoard: %v", err)
	}

	parent := boards.New(ctx, svc)
	parent.List.SetItems(boards.NewList(&[]service.Board{*board}))
	parent.List.Select(0)

	menu := New(ctx, svc, &parent)
	menu.List.SetItems(NewList(items))
	menu.List.Select(0)

	var captured []byte
	prevWrite := writeClipboardText
	writeClipboardText = func(data []byte) { captured = append([]byte{}, data...) }
	defer func() { writeClipboardText = prevWrite }()

	cmd := menu.Copy()
	if cmd != nil {
		cmd()
	}

	if len(captured) == 0 {
		t.Fatalf("expected clipboard content to be written")
	}

	var saved service.Item
	if err := json.Unmarshal(captured, &saved); err != nil {
		t.Fatalf("Unmarshal copied item: %v", err)
	}

	if saved.Title != item.Title || saved.Description != item.Description {
		t.Fatalf("unexpected copied item: %+v", saved)
	}
	if !saved.Completed {
		t.Fatalf("expected completed item in clipboard")
	}
	if len(saved.Tags) != len(item.Tags) {
		t.Fatalf("expected tags copied, got %+v", saved.Tags)
	}
}

func TestPasteCreatesItemOnSelectedBoard(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	ctx := testutil.MustContext()
	board, err := svc.CreateBoard(ctx, "Inbox")
	if err != nil {
		t.Fatalf("CreateBoard: %v", err)
	}

	parent := boards.New(ctx, svc)
	parent.List.SetItems(boards.NewList(&[]service.Board{*board}))
	parent.List.Select(0)

	menu := New(ctx, svc, &parent)

	clipItem := service.Item{
		Item: repository.Item{
			Title:       "from clipboard",
			Description: "pasted description",
			Completed:   true,
		},
		Tags: []string{"work", "go"},
	}
	data, err := json.Marshal(clipItem)
	if err != nil {
		t.Fatalf("Marshal clipboard item: %v", err)
	}

	prevRead := readClipboardText
	readClipboardText = func() []byte { return data }
	defer func() { readClipboardText = prevRead }()

	cmd := menu.Paste()
	if cmd == nil {
		t.Fatalf("expected paste command")
	}

	msg := cmd()
	created, ok := msg.(CreateItemMsg)
	if !ok {
		t.Fatalf("expected CreateItemMsg, got %T", msg)
	}
	if created.Error != nil {
		t.Fatalf("Paste returned error: %v", created.Error)
	}
	if created.Item == nil {
		t.Fatalf("expected created item")
	}
	if created.Item.BoardID != board.ID {
		t.Fatalf("expected item on board %d, got %d", board.ID, created.Item.BoardID)
	}
	if created.Item.Title != clipItem.Title || created.Item.Description != clipItem.Description {
		t.Fatalf("unexpected item fields: %+v", created.Item)
	}
	if !created.Item.Completed {
		t.Fatalf("expected completed item after paste")
	}
	if len(created.Item.Tags) != len(clipItem.Tags) {
		t.Fatalf("expected tags to persist, got %+v", created.Item.Tags)
	}

	items, err := svc.ListItemsByBoard(ctx, board)
	if err != nil {
		t.Fatalf("ListItemsByBoard: %v", err)
	}
	if len(*items) != 1 {
		t.Fatalf("expected 1 item in board after paste, got %d", len(*items))
	}
}

func TestItemCreateRenameToggleTagDeleteFlow(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	ctx := testutil.MustContext()
	board, err := svc.CreateBoard(ctx, "Work")
	if err != nil {
		t.Fatalf("CreateBoard: %v", err)
	}

	parent := boards.New(ctx, svc)
	parent.List.SetItems(boards.NewList(&[]service.Board{*board}))
	parent.List.Select(0)

	menu := New(ctx, svc, &parent)

	// Create item via input flow.
	menu.InitCreateItem()
	menu.Input.SetValue("Title")
	model, cmd := menu.Update(tea.KeyMsg{Type: tea.KeyEnter})
	menu = model.(MenuModel)
	if menu.Context.State != CreateItemDescState {
		t.Fatalf("expected to prompt for description, got state %v", menu.Context.State)
	}
	menu.Input.SetValue("Desc")
	model, cmd = menu.Update(tea.KeyMsg{Type: tea.KeyEnter})
	menu = model.(MenuModel)
	if cmd == nil {
		t.Fatalf("expected create command")
	}
	if msg := cmd(); msg != nil {
		model, _ = menu.Update(msg)
		menu = model.(MenuModel)
	}
	if len(menu.List.Items()) != 1 {
		t.Fatalf("expected 1 item after create, got %d", len(menu.List.Items()))
	}

	// Rename item via two-step input.
	menu.InitRenameItem()
	menu.Input.SetValue("New Title")
	model, cmd = menu.Update(tea.KeyMsg{Type: tea.KeyEnter})
	menu = model.(MenuModel)
	if menu.Context.State != RenameItemDescState {
		t.Fatalf("expected rename description state, got %v", menu.Context.State)
	}
	menu.Input.SetValue("New Desc")
	model, cmd = menu.Update(tea.KeyMsg{Type: tea.KeyEnter})
	menu = model.(MenuModel)
	if cmd == nil {
		t.Fatalf("expected rename command")
	}
	if msg := cmd(); msg != nil {
		model, _ = menu.Update(msg)
		menu = model.(MenuModel)
	}
	renamed, _ := menu.selectedItem()
	if renamed.Itm.Title != "New Title" || renamed.Itm.Description != "New Desc" {
		t.Fatalf("rename not applied, got %+v", renamed.Itm)
	}

	// Toggle completion.
	cmd = menu.ToggleComplete()
	if cmd == nil {
		t.Fatalf("expected toggle command")
	}
	if msg := cmd(); msg != nil {
		model, _ = menu.Update(msg)
		menu = model.(MenuModel)
	}
	toggled, _ := menu.selectedItem()
	if !toggled.Itm.Completed {
		t.Fatalf("expected item to be completed after toggle")
	}

	// Update tags.
	menu.InitUpdateTags()
	menu.Input.SetValue("one, two")
	model, cmd = menu.Update(tea.KeyMsg{Type: tea.KeyEnter})
	menu = model.(MenuModel)
	if cmd == nil {
		t.Fatalf("expected update tags command")
	}
	if msg := cmd(); msg != nil {
		model, _ = menu.Update(msg)
		menu = model.(MenuModel)
	}
	tagged, _ := menu.selectedItem()
	if len(tagged.Itm.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %+v", tagged.Itm.Tags)
	}

	// Delete item (clipboard write stubbed).
	prevWrite := writeClipboardText
	writeClipboardText = func(data []byte) {}
	defer func() { writeClipboardText = prevWrite }()

	cmd = menu.DeleteItem()
	msg := cmd()
	model, _ = menu.Update(msg)
	menu = model.(MenuModel)

	refresh := menu.ListItems()
	msg = refresh()
	model, _ = menu.Update(msg)
	menu = model.(MenuModel)
	if len(menu.List.Items()) != 0 {
		t.Fatalf("expected list empty after delete, got %d", len(menu.List.Items()))
	}
	items, err := svc.ListItemsByBoard(ctx, board)
	if err != nil {
		t.Fatalf("ListItemsByBoard: %v", err)
	}
	if len(*items) != 0 {
		t.Fatalf("expected no items in service after delete, got %d", len(*items))
	}
}
