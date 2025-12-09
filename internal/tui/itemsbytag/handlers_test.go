package itemsbytag

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/testutil"
	"github.com/rhajizada/donezo/internal/tui/navigation"
	"github.com/rhajizada/donezo/internal/tui/tags"
)

func newItemsByTagMenu(t *testing.T) (MenuModel, func()) {
	t.Helper()
	svc, cleanup := testutil.NewTestService(t)
	ctx := testutil.MustContext()
	var err error

	board, err := svc.CreateBoard(ctx, "Inbox")
	if err != nil {
		t.Fatalf("CreateBoard: %v", err)
	}
	item, err := svc.CreateItem(ctx, board, "task", "desc")
	if err != nil {
		t.Fatalf("CreateItem: %v", err)
	}
	item.Tags = []string{"work"}
	_, err = svc.UpdateItem(ctx, item)
	if err != nil {
		t.Fatalf("UpdateItem: %v", err)
	}

	tagCount, _ := svc.CountItemsByTag(ctx, "work")
	parent := tags.NewModel(ctx, svc)
	parent.List.SetItems(tags.NewList([]tags.Item{tags.NewItem("work", tagCount)}))
	parent.List.Select(0)

	menu := New(ctx, svc, &parent)
	menu.List.SetItems(NewList(&[]service.Item{*item}))
	menu.List.Select(0)
	return menu, cleanup
}

//nolint:gocognit // covering keybinding branches
func TestItemsByTagKeyBindings(t *testing.T) {
	t.Run("back sends BackMsg", func(t *testing.T) {
		menu, cleanup := newItemsByTagMenu(t)
		defer cleanup()

		_, cmd := menu.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		if cmd == nil {
			t.Fatalf("expected command")
		}
		if msg := cmd(); msg != (navigation.BackMsg{}) {
			t.Fatalf("expected BackMsg, got %v", msg)
		}
	})

	t.Run("rename and update tags enter states", func(t *testing.T) {
		menu, cleanup := newItemsByTagMenu(t)
		defer cleanup()

		model, _ := menu.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
		menu = model.(MenuModel)
		if menu.Context.State != RenameItemNameState {
			t.Fatalf("expected rename state, got %v", menu.Context.State)
		}

		menu.Context.State = DefaultState
		model, _ = menu.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
		menu = model.(MenuModel)
		if menu.Context.State != UpdateTagsState {
			t.Fatalf("expected update tags state, got %v", menu.Context.State)
		}
	})

	t.Run("delete, refresh, toggle, navigation", func(t *testing.T) {
		menu, cleanup := newItemsByTagMenu(t)
		defer cleanup()

		_, cmd := menu.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
		if cmd == nil {
			t.Fatalf("expected delete cmd")
		}
		if msg := cmd(); msg == nil {
			t.Fatalf("expected DeleteItemMsg")
		} else if _, ok := msg.(DeleteItemMsg); !ok {
			t.Fatalf("expected DeleteItemMsg, got %T", msg)
		}

		_, cmd = menu.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'R'}})
		if cmd == nil {
			t.Fatalf("expected refresh cmd")
		}
		if msg := cmd(); msg == nil {
			t.Fatalf("expected ListItemsMsg")
		} else if _, ok := msg.(ListItemsMsg); !ok {
			t.Fatalf("expected ListItemsMsg, got %T", msg)
		}

		_, cmd = menu.Update(tea.KeyMsg{Type: tea.KeySpace})
		if cmd == nil {
			t.Fatalf("expected toggle cmd")
		}
		if msg := cmd(); msg == nil {
			t.Fatalf("expected ToggleItemMsg")
		} else if _, ok := msg.(ToggleItemMsg); !ok {
			t.Fatalf("expected ToggleItemMsg, got %T", msg)
		}

		_, cmd = menu.Update(tea.KeyMsg{Type: tea.KeyTab})
		if cmd == nil {
			t.Fatalf("expected next tag cmd")
		}
		if msg := cmd(); msg != (navigation.TagDeltaMsg{Delta: 1}) {
			t.Fatalf("expected TagDeltaMsg +1, got %v", msg)
		}

		_, cmd = menu.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
		if cmd == nil {
			t.Fatalf("expected prev tag cmd")
		}
		if msg := cmd(); msg != (navigation.TagDeltaMsg{Delta: -1}) {
			t.Fatalf("expected TagDeltaMsg -1, got %v", msg)
		}
	})
}
