package tags

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rhajizada/donezo/internal/testutil"
	"github.com/rhajizada/donezo/internal/tui/navigation"
)

func newTagMenu(t *testing.T) (MenuModel, func()) {
	t.Helper()
	svc, cleanup := testutil.NewTestService(t)
	ctx := testutil.MustContext()
	var err error

	board, err := svc.CreateBoard(ctx, "Inbox")
	require.NoError(t, err)
	item, err := svc.CreateItem(ctx, board, "task", "desc")
	require.NoError(t, err)
	item.Tags = []string{"work"}
	_, err = svc.UpdateItem(ctx, item)
	require.NoError(t, err)

	tagCount, _ := svc.CountItemsByTag(ctx, "work")

	menu := NewModel(ctx, svc)
	menu.List.SetItems(NewList([]Item{NewItem("work", tagCount)}))
	menu.List.Select(0)
	return menu, cleanup
}

func TestTagsKeyBindings(t *testing.T) {
	tests := []struct {
		name      string
		msg       tea.KeyPressMsg
		assertCmd func(*testing.T, tea.Cmd)
	}{
		{
			name: "enter opens items",
			msg:  tea.KeyPressMsg{Code: tea.KeyEnter},
			assertCmd: func(t *testing.T, cmd tea.Cmd) {
				require.NotNil(t, cmd)
				assert.Equal(t, navigation.OpenTagItemsMsg{}, cmd())
			},
		},
		{
			name: "tab switches to boards",
			msg:  tea.KeyPressMsg{Code: tea.KeyTab},
			assertCmd: func(t *testing.T, cmd tea.Cmd) {
				require.NotNil(t, cmd)
				assert.Equal(t, navigation.SwitchMainViewMsg{View: navigation.ViewBoards}, cmd())
			},
		},
		{
			name: "delete sends delete tag message",
			msg:  tea.KeyPressMsg{Code: 'd', Text: "d"},
			assertCmd: func(t *testing.T, cmd tea.Cmd) {
				require.NotNil(t, cmd)
				_, ok := cmd().(DeleteTagMsg)
				assert.True(t, ok)
			},
		},
		{
			name: "refresh sends list tags message",
			msg:  tea.KeyPressMsg{Code: 'R', Text: "R"},
			assertCmd: func(t *testing.T, cmd tea.Cmd) {
				require.NotNil(t, cmd)
				_, ok := cmd().(ListTagsMsg)
				assert.True(t, ok)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menu, cleanup := newTagMenu(t)
			defer cleanup()

			_, cmd := menu.Update(tt.msg)
			tt.assertCmd(t, cmd)
		})
	}

	t.Run("copy writes to clipboard", func(t *testing.T) {
		menu, cleanup := newTagMenu(t)
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
