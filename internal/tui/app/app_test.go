package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/testutil"
	"github.com/rhajizada/donezo/internal/tui/boards"
	"github.com/rhajizada/donezo/internal/tui/itemsbytag"
	"github.com/rhajizada/donezo/internal/tui/navigation"
	"github.com/rhajizada/donezo/internal/tui/tags"
)

func seedBoard(t *testing.T, svc *service.Service, name string) *service.Board {
	t.Helper()
	board, err := svc.CreateBoard(testutil.MustContext(), name)
	if err != nil {
		t.Fatalf("CreateBoard: %v", err)
	}
	return board
}

func TestAppNavigatesBetweenBoardsAndItems(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	board := seedBoard(t, svc, "Inbox")
	if _, err := svc.CreateItem(testutil.MustContext(), board, "task", "desc"); err != nil {
		t.Fatalf("CreateItem: %v", err)
	}

	ctx := testutil.MustContext()
	m := New(ctx, svc)

	// Preload boards list so selection exists.
	m.boards.List.SetItems(boards.NewList(&[]service.Board{*board}))

	model, _ := m.Update(navigation.OpenBoardItemsMsg{})
	am, ok := model.(AppModel)
	if !ok {
		t.Fatalf("unexpected model type %T", model)
	}
	if am.active != navigation.ViewItemsByBoard {
		t.Fatalf("expected active view items-by-board, got %v", am.active)
	}
	if am.itemsByBoard == nil {
		t.Fatalf("expected itemsByBoard to be initialized")
	}

	// Navigate back to boards view.
	model, _ = am.Update(navigation.BackMsg{})
	am, ok = model.(AppModel)
	if !ok {
		t.Fatalf("unexpected model type after back %T", model)
	}
	if am.active != navigation.ViewBoards {
		t.Fatalf("expected active view boards, got %v", am.active)
	}
}

func TestBoardsEscDoesNotQuit(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	board := seedBoard(t, svc, "Inbox")
	ctx := testutil.MustContext()
	menu := boards.New(ctx, svc)
	menu.List.SetItems(boards.NewList(&[]service.Board{*board}))

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	model, cmd := menu.Update(msg)
	if cmd != nil {
		// ESC on boards view should not emit a quit command.
		if _, ok := cmd().(tea.QuitMsg); ok {
			t.Fatalf("esc should not quit application")
		}
	}

	if updated, ok := model.(boards.MenuModel); ok {
		if updated.State != boards.DefaultState {
			t.Fatalf("expected default state after esc, got %v", updated.State)
		}
	} else {
		t.Fatalf("unexpected model type %T", model)
	}
}

func TestTagsEscDoesNotQuit(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	board := seedBoard(t, svc, "Inbox")
	item, err := svc.CreateItem(testutil.MustContext(), board, "task", "desc")
	if err != nil {
		t.Fatalf("CreateItem: %v", err)
	}
	item.Tags = []string{"inbox"}
	if _, err := svc.UpdateItem(testutil.MustContext(), item); err != nil {
		t.Fatalf("UpdateItem (add tag): %v", err)
	}

	tagCount, err := svc.CountItemsByTag(testutil.MustContext(), "inbox")
	if err != nil {
		t.Fatalf("CountItemsByTag: %v", err)
	}

	ctx := testutil.MustContext()
	menu := tags.NewModel(ctx, svc)
	menu.List.SetItems(tags.NewList([]tags.Item{
		tags.NewItem("inbox", tagCount),
	}))

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	model, cmd := menu.Update(msg)
	if cmd != nil {
		if _, ok := cmd().(tea.QuitMsg); ok {
			t.Fatalf("esc should not quit application in tags view")
		}
	}

	if _, ok := model.(tags.MenuModel); !ok {
		t.Fatalf("unexpected model type %T", model)
	}
}

func TestBoardsTabSwitchesToTags(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	ctx := testutil.MustContext()
	board := seedBoard(t, svc, "Inbox")
	item, err := svc.CreateItem(ctx, board, "task", "desc")
	if err != nil {
		t.Fatalf("CreateItem: %v", err)
	}
	item.Tags = []string{"work"}
	if _, err := svc.UpdateItem(ctx, item); err != nil {
		t.Fatalf("UpdateItem: %v", err)
	}

	m := New(ctx, svc)
	m.boards.List.SetItems(boards.NewList(&[]service.Board{*board}))

	model, cmd := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	appModel := model.(AppModel)
	if cmd == nil {
		t.Fatalf("expected command from tab")
	}

	switchMsg := cmd()
	model, cmd = appModel.Update(switchMsg)
	appModel = model.(AppModel)
	if cmd != nil {
		if msg := cmd(); msg != nil {
			model, _ = appModel.Update(msg)
			appModel = model.(AppModel)
		}
	}

	if appModel.active != navigation.ViewTags {
		t.Fatalf("expected to switch to tags, got %v", appModel.active)
	}
	if len(appModel.tags.List.Items()) != 1 {
		t.Fatalf("expected tags to be loaded, got %d", len(appModel.tags.List.Items()))
	}
}

func TestItemsByBoardTabAndShiftTabCycleBoards(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	ctx := testutil.MustContext()
	board1 := seedBoard(t, svc, "One")
	board2 := seedBoard(t, svc, "Two")
	if _, err := svc.CreateItem(ctx, board1, "a", ""); err != nil {
		t.Fatalf("CreateItem board1: %v", err)
	}
	if _, err := svc.CreateItem(ctx, board2, "b", ""); err != nil {
		t.Fatalf("CreateItem board2: %v", err)
	}

	m := New(ctx, svc)
	m.boards.List.SetItems(boards.NewList(&[]service.Board{*board1, *board2}))
	model, _ := m.Update(navigation.OpenBoardItemsMsg{})
	appModel := model.(AppModel)

	// Move forward (tab)
	model, cmd := appModel.Update(navigation.BoardDeltaMsg{Delta: 1})
	appModel = model.(AppModel)
	if cmd != nil {
		if msg := cmd(); msg != nil {
			model, _ = appModel.Update(msg)
			appModel = model.(AppModel)
		}
	}
	if appModel.boards.List.Index() != 1 {
		t.Fatalf("expected board index 1, got %d", appModel.boards.List.Index())
	}
	if appModel.itemsByBoard.List.Title != board2.Name {
		t.Fatalf("expected items view for board %s, got %s", board2.Name, appModel.itemsByBoard.List.Title)
	}

	// Move backward (shift+tab)
	model, cmd = appModel.Update(navigation.BoardDeltaMsg{Delta: -1})
	appModel = model.(AppModel)
	if cmd != nil {
		if msg := cmd(); msg != nil {
			model, _ = appModel.Update(msg)
			appModel = model.(AppModel)
		}
	}
	if appModel.boards.List.Index() != 0 {
		t.Fatalf("expected board index 0, got %d", appModel.boards.List.Index())
	}
	if appModel.itemsByBoard.List.Title != board1.Name {
		t.Fatalf("expected items view for board %s after reverse, got %s", board1.Name, appModel.itemsByBoard.List.Title)
	}
}

func TestItemsByTagNextAndPrevious(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	ctx := testutil.MustContext()
	board := seedBoard(t, svc, "Inbox")
	itemA, err := svc.CreateItem(ctx, board, "a", "")
	if err != nil {
		t.Fatalf("CreateItem a: %v", err)
	}
	itemA.Tags = []string{"alpha"}
	if _, err := svc.UpdateItem(ctx, itemA); err != nil {
		t.Fatalf("UpdateItem a: %v", err)
	}
	itemB, err := svc.CreateItem(ctx, board, "b", "")
	if err != nil {
		t.Fatalf("CreateItem b: %v", err)
	}
	itemB.Tags = []string{"beta"}
	if _, err := svc.UpdateItem(ctx, itemB); err != nil {
		t.Fatalf("UpdateItem b: %v", err)
	}

	m := New(ctx, svc)

	tagCountAlpha, _ := svc.CountItemsByTag(ctx, "alpha")
	tagCountBeta, _ := svc.CountItemsByTag(ctx, "beta")
	m.tags.List.SetItems(tags.NewList([]tags.Item{
		tags.NewItem("alpha", tagCountAlpha),
		tags.NewItem("beta", tagCountBeta),
	}))

	model, _ := m.Update(navigation.OpenTagItemsMsg{})
	appModel := model.(AppModel)
	appModel.tags.List.Select(0)
	itemTagModel := itemsbytag.New(ctx, svc, appModel.tags)
	appModel.itemsByTag = &itemTagModel
	appModel.active = navigation.ViewItemsByTag

	// Move to next tag
	model, cmd := appModel.Update(navigation.TagDeltaMsg{Delta: 1})
	appModel = model.(AppModel)
	if cmd != nil {
		if msg := cmd(); msg != nil {
			model, _ = appModel.Update(msg)
			appModel = model.(AppModel)
		}
	}
	if appModel.tags.List.Index() != 1 {
		t.Fatalf("expected tag index 1, got %d", appModel.tags.List.Index())
	}
	if appModel.itemsByTag.List.Title != "beta" {
		t.Fatalf("expected items for beta tag, got %s", appModel.itemsByTag.List.Title)
	}

	// Move back to previous tag
	model, cmd = appModel.Update(navigation.TagDeltaMsg{Delta: -1})
	appModel = model.(AppModel)
	if cmd != nil {
		if msg := cmd(); msg != nil {
			model, _ = appModel.Update(msg)
			appModel = model.(AppModel)
		}
	}
	if appModel.tags.List.Index() != 0 {
		t.Fatalf("expected tag index 0, got %d", appModel.tags.List.Index())
	}
	if appModel.itemsByTag.List.Title != "alpha" {
		t.Fatalf("expected items for alpha tag, got %s", appModel.itemsByTag.List.Title)
	}
}
