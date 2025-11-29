package boards

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/testutil"
	"github.com/rhajizada/donezo/internal/tui/navigation"
)

func newBoardMenu(t *testing.T) (MenuModel, func()) {
	t.Helper()
	svc, cleanup := testutil.NewTestService(t)
	ctx := testutil.MustContext()
	board, err := svc.CreateBoard(ctx, "Inbox")
	if err != nil {
		t.Fatalf("CreateBoard: %v", err)
	}
	menu := New(ctx, svc)
	menu.List.SetItems(NewList(&[]service.Board{*board}))
	menu.List.Select(0)
	return menu, cleanup
}

//nolint:gocognit // covering multiple keybinding branches
func TestBoardsKeyBindings(t *testing.T) {
	t.Run("enter opens items", func(t *testing.T) {
		menu, cleanup := newBoardMenu(t)
		defer cleanup()

		model, cmd := menu.Update(tea.KeyMsg{Type: tea.KeyEnter})
		if cmd == nil {
			t.Fatalf("expected command for enter")
		}
		menu = model.(MenuModel)
		msg := cmd()
		if _, ok := msg.(navigation.OpenBoardItemsMsg); !ok {
			t.Fatalf("expected OpenBoardItemsMsg, got %T", msg)
		}
		if menu.State != DefaultState {
			t.Fatalf("expected default state, got %v", menu.State)
		}
	})

	t.Run("tab switches to tags", func(t *testing.T) {
		menu, cleanup := newBoardMenu(t)
		defer cleanup()

		_, cmd := menu.Update(tea.KeyMsg{Type: tea.KeyTab})
		if cmd == nil {
			t.Fatalf("expected command for tab")
		}
		if msg := cmd(); msg != (navigation.SwitchMainViewMsg{View: navigation.ViewTags}) {
			t.Fatalf("expected SwitchMainViewMsg to tags, got %v", msg)
		}
	})

	t.Run("create and rename set state", func(t *testing.T) {
		menu, cleanup := newBoardMenu(t)
		defer cleanup()

		model, _ := menu.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
		menu = model.(MenuModel)
		if menu.State != CreateBoardState {
			t.Fatalf("expected create state, got %v", menu.State)
		}

		menu.State = DefaultState
		model, _ = menu.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
		menu = model.(MenuModel)
		if menu.State != RenameBoardState {
			t.Fatalf("expected rename state, got %v", menu.State)
		}
	})

	t.Run("delete sends DeleteBoardMsg", func(t *testing.T) {
		menu, cleanup := newBoardMenu(t)
		defer cleanup()

		_, cmd := menu.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
		if cmd == nil {
			t.Fatalf("expected delete command")
		}
		if msg := cmd(); msg == nil {
			t.Fatalf("expected DeleteBoardMsg")
		} else if _, ok := msg.(DeleteBoardMsg); !ok {
			t.Fatalf("expected DeleteBoardMsg, got %T", msg)
		}
	})

	t.Run("refresh sends ListBoardsMsg", func(t *testing.T) {
		menu, cleanup := newBoardMenu(t)
		defer cleanup()

		_, cmd := menu.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'R'}})
		if cmd == nil {
			t.Fatalf("expected refresh command")
		}
		if msg := cmd(); msg == nil {
			t.Fatalf("expected ListBoardsMsg")
		} else if _, ok := msg.(ListBoardsMsg); !ok {
			t.Fatalf("expected ListBoardsMsg, got %T", msg)
		}
	})

	t.Run("copy writes to clipboard", func(t *testing.T) {
		menu, cleanup := newBoardMenu(t)
		defer cleanup()

		var captured []byte
		prev := writeClipboardText
		writeClipboardText = func(data []byte) { captured = append([]byte{}, data...) }
		defer func() { writeClipboardText = prev }()

		_, cmd := menu.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
		if cmd == nil {
			t.Fatalf("expected copy command")
		}
		cmd()
		if len(captured) == 0 {
			t.Fatalf("expected clipboard write")
		}
	})
}
