package tags

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/rhajizada/donezo/internal/testutil"
	"github.com/rhajizada/donezo/internal/tui/navigation"
)

func newTagMenu(t *testing.T) (MenuModel, func()) {
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

	menu := NewModel(ctx, svc)
	menu.List.SetItems(NewList([]Item{NewItem("work", tagCount)}))
	menu.List.Select(0)
	return menu, cleanup
}

//nolint:gocognit // covering keybinding branches
func TestTagsKeyBindings(t *testing.T) {
	t.Run("enter opens items", func(t *testing.T) {
		menu, cleanup := newTagMenu(t)
		defer cleanup()

		_, cmd := menu.Update(tea.KeyMsg{Type: tea.KeyEnter})
		if cmd == nil {
			t.Fatalf("expected command")
		}
		if msg := cmd(); msg != (navigation.OpenTagItemsMsg{}) {
			t.Fatalf("expected OpenTagItemsMsg, got %v", msg)
		}
	})

	t.Run("tab switches to boards", func(t *testing.T) {
		menu, cleanup := newTagMenu(t)
		defer cleanup()

		_, cmd := menu.Update(tea.KeyMsg{Type: tea.KeyTab})
		if cmd == nil {
			t.Fatalf("expected command")
		}
		if msg := cmd(); msg != (navigation.SwitchMainViewMsg{View: navigation.ViewBoards}) {
			t.Fatalf("expected switch to boards, got %v", msg)
		}
	})

	t.Run("delete and refresh and copy", func(t *testing.T) {
		menu, cleanup := newTagMenu(t)
		defer cleanup()

		_, cmd := menu.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
		if cmd == nil {
			t.Fatalf("expected delete command")
		}
		if msg := cmd(); msg == nil {
			t.Fatalf("expected DeleteTagMsg")
		} else if _, ok := msg.(DeleteTagMsg); !ok {
			t.Fatalf("expected DeleteTagMsg, got %T", msg)
		}

		_, cmd = menu.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'R'}})
		if cmd == nil {
			t.Fatalf("expected refresh command")
		}
		if msg := cmd(); msg == nil {
			t.Fatalf("expected ListTagsMsg")
		} else if _, ok := msg.(ListTagsMsg); !ok {
			t.Fatalf("expected ListTagsMsg, got %T", msg)
		}

		var captured []byte
		prev := writeClipboardText
		writeClipboardText = func(data []byte) { captured = append([]byte{}, data...) }
		defer func() { writeClipboardText = prev }()

		_, cmd = menu.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
		if cmd == nil {
			t.Fatalf("expected copy command")
		}
		cmd()
		if len(captured) == 0 {
			t.Fatalf("expected clipboard write")
		}
	})
}
