package itemsbyboard

import (
	"encoding/json"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rhajizada/donezo/internal/repository"
	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/testutil"
	"github.com/rhajizada/donezo/internal/tui/boards"
	"github.com/rhajizada/donezo/internal/tui/navigation"
)

func newItemMenu(t *testing.T) (MenuModel, func()) {
	t.Helper()
	svc, cleanup := testutil.NewTestService(t)
	ctx := testutil.MustContext()
	board, err := svc.CreateBoard(ctx, "Inbox")
	require.NoError(t, err)
	item, err := svc.CreateItem(ctx, board, "task", "desc")
	require.NoError(t, err)

	parent := boards.New(ctx, svc)
	parent.List.SetItems(boards.NewList(&[]service.Board{*board}))
	parent.List.Select(0)

	menu := New(ctx, svc, &parent)
	menu.List.SetItems(NewList(&[]service.Item{*item}))
	menu.List.Select(0)
	return menu, cleanup
}

func TestItemsByBoardKeyBindings(t *testing.T) {
	tests := []struct {
		name        string
		msg         tea.KeyPressMsg
		assertModel func(*testing.T, MenuModel)
		assertCmd   func(*testing.T, tea.Cmd)
	}{
		{
			name: "back sends back message",
			msg:  tea.KeyPressMsg{Code: tea.KeyBackspace},
			assertCmd: func(t *testing.T, cmd tea.Cmd) {
				require.NotNil(t, cmd)
				assert.Equal(t, navigation.BackMsg{}, cmd())
			},
		},
		{
			name: "create enters create state",
			msg:  tea.KeyPressMsg{Code: 'a', Text: "a"},
			assertModel: func(t *testing.T, menu MenuModel) {
				assert.Equal(t, CreateItemNameState, menu.Context.State)
			},
		},
		{
			name: "rename enters rename state",
			msg:  tea.KeyPressMsg{Code: 'r', Text: "r"},
			assertModel: func(t *testing.T, menu MenuModel) {
				assert.Equal(t, RenameItemNameState, menu.Context.State)
			},
		},
		{
			name: "update tags enters update tags state",
			msg:  tea.KeyPressMsg{Code: 't', Text: "t"},
			assertModel: func(t *testing.T, menu MenuModel) {
				assert.Equal(t, UpdateTagsState, menu.Context.State)
			},
		},
		{
			name: "delete sends delete message",
			msg:  tea.KeyPressMsg{Code: 'd', Text: "d"},
			assertCmd: func(t *testing.T, cmd tea.Cmd) {
				require.NotNil(t, cmd)
				_, ok := cmd().(DeleteItemMsg)
				assert.True(t, ok)
			},
		},
		{
			name: "refresh sends list items message",
			msg:  tea.KeyPressMsg{Code: 'R', Text: "R"},
			assertCmd: func(t *testing.T, cmd tea.Cmd) {
				require.NotNil(t, cmd)
				_, ok := cmd().(ListItemsMsg)
				assert.True(t, ok)
			},
		},
		{
			name: "space toggles completion",
			msg:  tea.KeyPressMsg{Code: tea.KeySpace, Text: " "},
			assertCmd: func(t *testing.T, cmd tea.Cmd) {
				require.NotNil(t, cmd)
				_, ok := cmd().(ToggleItemMsg)
				assert.True(t, ok)
			},
		},
		{
			name: "tab moves to next board",
			msg:  tea.KeyPressMsg{Code: tea.KeyTab},
			assertCmd: func(t *testing.T, cmd tea.Cmd) {
				require.NotNil(t, cmd)
				assert.Equal(t, navigation.BoardDeltaMsg{Delta: 1}, cmd())
			},
		},
		{
			name: "shift tab moves to previous board",
			msg:  tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift},
			assertCmd: func(t *testing.T, cmd tea.Cmd) {
				require.NotNil(t, cmd)
				assert.Equal(t, navigation.BoardDeltaMsg{Delta: -1}, cmd())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menu, cleanup := newItemMenu(t)
			defer cleanup()

			if tt.name == "update tags enters update tags state" {
				menu.Context.State = DefaultState
			}

			model, cmd := menu.Update(tt.msg)
			menu = model.(MenuModel)
			if tt.assertModel != nil {
				tt.assertModel(t, menu)
			}
			if tt.assertCmd != nil {
				tt.assertCmd(t, cmd)
			}
		})
	}

	t.Run("copy and paste", func(t *testing.T) {
		menu, cleanup := newItemMenu(t)
		defer cleanup()

		var captured []byte
		prevWrite := writeClipboardText
		writeClipboardText = func(data []byte) { captured = append([]byte{}, data...) }
		defer func() { writeClipboardText = prevWrite }()

		_, cmd := menu.Update(tea.KeyPressMsg{Code: 'y', Text: "y"})
		require.NotNil(t, cmd)
		cmd()
		assert.NotEmpty(t, captured)

		clipItem := service.Item{
			Item: repository.Item{
				Title:       "pasted",
				Description: "desc",
			},
		}
		data, err := json.Marshal(clipItem)
		require.NoError(t, err)
		prevRead := readClipboardText
		readClipboardText = func() []byte { return data }
		defer func() { readClipboardText = prevRead }()

		_, cmd = menu.Update(tea.KeyPressMsg{Code: 'p', Text: "p"})
		require.NotNil(t, cmd)
		_, ok := cmd().(CreateItemMsg)
		assert.True(t, ok)
	})
}
