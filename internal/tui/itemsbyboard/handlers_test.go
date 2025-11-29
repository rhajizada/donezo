package itemsbyboard

import (
	"encoding/json"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/rhajizada/donezo/internal/repository"
	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/testutil"
	"github.com/rhajizada/donezo/internal/tui/boards"
	"github.com/rhajizada/donezo/internal/tui/navigation"
)

func newItemMenu(t *testing.T) (MenuModel, *service.Service, func()) {
	t.Helper()
	svc, cleanup := testutil.NewTestService(t)
	ctx := testutil.MustContext()
	board, err := svc.CreateBoard(ctx, "Inbox")
	if err != nil {
		t.Fatalf("CreateBoard: %v", err)
	}
	item, err := svc.CreateItem(ctx, board, "task", "desc")
	if err != nil {
		t.Fatalf("CreateItem: %v", err)
	}

	parent := boards.New(ctx, svc)
	parent.List.SetItems(boards.NewList(&[]service.Board{*board}))
	parent.List.Select(0)

	menu := New(ctx, svc, &parent)
	menu.List.SetItems(NewList(&[]service.Item{*item}))
	menu.List.Select(0)
	return menu, svc, cleanup
}

func TestItemsByBoardKeyBindings(t *testing.T) {
	t.Run("back sends BackMsg", func(t *testing.T) {
		menu, _, cleanup := newItemMenu(t)
		defer cleanup()

		_, cmd := menu.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		if cmd == nil {
			t.Fatalf("expected command")
		}
		if msg := cmd(); msg != (navigation.BackMsg{}) {
			t.Fatalf("expected BackMsg, got %v", msg)
		}
	})

	t.Run("create enters input states", func(t *testing.T) {
		menu, _, cleanup := newItemMenu(t)
		defer cleanup()

		model, _ := menu.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
		menu = model.(MenuModel)
		if menu.Context.State != CreateItemNameState {
			t.Fatalf("expected create name state, got %v", menu.Context.State)
		}
	})

	t.Run("rename and update tags enter states", func(t *testing.T) {
		menu, _, cleanup := newItemMenu(t)
		defer cleanup()

		model, _ := menu.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
		menu = model.(MenuModel)
		if menu.Context.State != RenameItemNameState {
			t.Fatalf("expected rename name state, got %v", menu.Context.State)
		}

		menu.Context.State = DefaultState
		model, _ = menu.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
		menu = model.(MenuModel)
		if menu.Context.State != UpdateTagsState {
			t.Fatalf("expected update tags state, got %v", menu.Context.State)
		}
	})

	t.Run("delete and refresh commands", func(t *testing.T) {
		menu, _, cleanup := newItemMenu(t)
		defer cleanup()

		_, cmd := menu.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
		if cmd == nil {
			t.Fatalf("expected delete command")
		}
		if msg := cmd(); msg == nil {
			t.Fatalf("expected DeleteItemMsg")
		} else if _, ok := msg.(DeleteItemMsg); !ok {
			t.Fatalf("expected DeleteItemMsg, got %T", msg)
		}

		_, cmd = menu.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'R'}})
		if cmd == nil {
			t.Fatalf("expected refresh command")
		}
		if msg := cmd(); msg == nil {
			t.Fatalf("expected ListItemsMsg")
		} else if _, ok := msg.(ListItemsMsg); !ok {
			t.Fatalf("expected ListItemsMsg, got %T", msg)
		}
	})

	t.Run("toggle complete and navigation", func(t *testing.T) {
		menu, _, cleanup := newItemMenu(t)
		defer cleanup()

		_, cmd := menu.Update(tea.KeyMsg{Type: tea.KeySpace})
		if cmd == nil {
			t.Fatalf("expected toggle command")
		}
		if msg := cmd(); msg == nil {
			t.Fatalf("expected ToggleItemMsg")
		} else if _, ok := msg.(ToggleItemMsg); !ok {
			t.Fatalf("expected ToggleItemMsg, got %T", msg)
		}

		_, cmd = menu.Update(tea.KeyMsg{Type: tea.KeyTab})
		if cmd == nil {
			t.Fatalf("expected next board command")
		}
		if msg := cmd(); msg != (navigation.BoardDeltaMsg{Delta: 1}) {
			t.Fatalf("expected BoardDeltaMsg +1, got %v", msg)
		}

		_, cmd = menu.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
		if cmd == nil {
			t.Fatalf("expected previous board command")
		}
		if msg := cmd(); msg != (navigation.BoardDeltaMsg{Delta: -1}) {
			t.Fatalf("expected BoardDeltaMsg -1, got %v", msg)
		}
	})

	t.Run("copy and paste", func(t *testing.T) {
		menu, _, cleanup := newItemMenu(t)
		defer cleanup()

		var captured []byte
		prevWrite := writeClipboardText
		writeClipboardText = func(data []byte) { captured = append([]byte{}, data...) }
		defer func() { writeClipboardText = prevWrite }()

		_, cmd := menu.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
		if cmd == nil {
			t.Fatalf("expected copy cmd")
		}
		cmd()
		if len(captured) == 0 {
			t.Fatalf("expected clipboard write")
		}

		clipItem := service.Item{
			Item: repository.Item{
				Title:       "pasted",
				Description: "desc",
			},
		}
		data, _ := json.Marshal(clipItem)
		prevRead := readClipboardText
		readClipboardText = func() []byte { return data }
		defer func() { readClipboardText = prevRead }()

		_, cmd = menu.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
		if cmd == nil {
			t.Fatalf("expected paste cmd")
		}
		if msg := cmd(); msg == nil {
			t.Fatalf("expected CreateItemMsg from paste")
		} else if _, ok := msg.(CreateItemMsg); !ok {
			t.Fatalf("expected CreateItemMsg, got %T", msg)
		}
	})
}
