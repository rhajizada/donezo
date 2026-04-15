package boards

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/testutil"
	"github.com/rhajizada/donezo/internal/tui/navigation"
	"github.com/rhajizada/donezo/internal/tui/styles"
)

func newBoardMenu(t *testing.T) (MenuModel, func()) {
	t.Helper()
	svc, cleanup := testutil.NewTestService(t)
	ctx := testutil.MustContext()
	board, err := svc.CreateBoard(ctx, "Inbox")
	require.NoError(t, err)
	menu := New(ctx, svc)
	menu.List.SetItems(NewList(&[]service.Board{*board}))
	menu.List.Select(0)
	return menu, cleanup
}

func TestBoardsKeyBindings(t *testing.T) {
	t.Run("window resize updates input width", func(t *testing.T) {
		menu, cleanup := newBoardMenu(t)
		defer cleanup()

		width := 120
		model, _ := menu.Update(tea.WindowSizeMsg{Width: width, Height: 40})
		menu = model.(MenuModel)

		h, _ := styles.App.GetFrameSize()
		assert.Equal(t, width-h, menu.Input.Width())
	})

	tests := []struct {
		name        string
		msg         tea.KeyPressMsg
		assertModel func(*testing.T, MenuModel)
		assertCmd   func(*testing.T, tea.Cmd)
	}{
		{
			name: "enter opens items",
			msg:  tea.KeyPressMsg{Code: tea.KeyEnter},
			assertModel: func(t *testing.T, menu MenuModel) {
				assert.Equal(t, DefaultState, menu.State)
			},
			assertCmd: func(t *testing.T, cmd tea.Cmd) {
				require.NotNil(t, cmd)
				_, ok := cmd().(navigation.OpenBoardItemsMsg)
				assert.True(t, ok)
			},
		},
		{
			name: "tab switches to tags",
			msg:  tea.KeyPressMsg{Code: tea.KeyTab},
			assertCmd: func(t *testing.T, cmd tea.Cmd) {
				require.NotNil(t, cmd)
				assert.Equal(t, navigation.SwitchMainViewMsg{View: navigation.ViewTags}, cmd())
			},
		},
		{
			name: "create enters create state",
			msg:  tea.KeyPressMsg{Code: 'a', Text: "a"},
			assertModel: func(t *testing.T, menu MenuModel) {
				assert.Equal(t, CreateBoardState, menu.State)
			},
		},
		{
			name: "rename enters rename state",
			msg:  tea.KeyPressMsg{Code: 'r', Text: "r"},
			assertModel: func(t *testing.T, menu MenuModel) {
				assert.Equal(t, RenameBoardState, menu.State)
			},
		},
		{
			name: "delete sends delete message",
			msg:  tea.KeyPressMsg{Code: 'd', Text: "d"},
			assertCmd: func(t *testing.T, cmd tea.Cmd) {
				require.NotNil(t, cmd)
				_, ok := cmd().(DeleteBoardMsg)
				assert.True(t, ok)
			},
		},
		{
			name: "refresh sends list message",
			msg:  tea.KeyPressMsg{Code: 'R', Text: "R"},
			assertCmd: func(t *testing.T, cmd tea.Cmd) {
				require.NotNil(t, cmd)
				_, ok := cmd().(ListBoardsMsg)
				assert.True(t, ok)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menu, cleanup := newBoardMenu(t)
			defer cleanup()

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

	t.Run("copy writes to clipboard", func(t *testing.T) {
		menu, cleanup := newBoardMenu(t)
		defer cleanup()

		var captured []byte
		prev := writeClipboardText
		writeClipboardText = func(data []byte) { captured = append([]byte{}, data...) }
		defer func() { writeClipboardText = prev }()

		_, cmd := menu.Update(tea.KeyPressMsg{Code: 'y', Text: "y"})
		require.NotNil(t, cmd)
		cmd()
		assert.NotEmpty(t, captured)
	})
}
