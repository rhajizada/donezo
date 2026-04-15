package itemsbytag

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	require.NoError(t, err)
	item, err := svc.CreateItem(ctx, board, "task", "desc")
	require.NoError(t, err)
	item.Tags = []string{"work"}
	_, err = svc.UpdateItem(ctx, item)
	require.NoError(t, err)

	tagCount, _ := svc.CountItemsByTag(ctx, "work")
	parent := tags.NewModel(ctx, svc)
	parent.List.SetItems(tags.NewList([]tags.Item{tags.NewItem("work", tagCount)}))
	parent.List.Select(0)

	menu := New(ctx, svc, &parent)
	menu.List.SetItems(NewList(&[]service.Item{*item}))
	menu.List.Select(0)
	return menu, cleanup
}

func TestItemsByTagKeyBindings(t *testing.T) {
	tests := []struct {
		name        string
		msg         tea.KeyPressMsg
		prepare     func(*MenuModel)
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
			name: "rename enters rename state",
			msg:  tea.KeyPressMsg{Code: 'r', Text: "r"},
			assertModel: func(t *testing.T, menu MenuModel) {
				assert.Equal(t, RenameItemNameState, menu.Context.State)
			},
		},
		{
			name: "update tags enters update tags state",
			msg:  tea.KeyPressMsg{Code: 't', Text: "t"},
			prepare: func(menu *MenuModel) {
				menu.Context.State = DefaultState
			},
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
			name: "space toggles item",
			msg:  tea.KeyPressMsg{Code: tea.KeySpace, Text: " "},
			assertCmd: func(t *testing.T, cmd tea.Cmd) {
				require.NotNil(t, cmd)
				_, ok := cmd().(ToggleItemMsg)
				assert.True(t, ok)
			},
		},
		{
			name: "tab moves to next tag",
			msg:  tea.KeyPressMsg{Code: tea.KeyTab},
			assertCmd: func(t *testing.T, cmd tea.Cmd) {
				require.NotNil(t, cmd)
				assert.Equal(t, navigation.TagDeltaMsg{Delta: 1}, cmd())
			},
		},
		{
			name: "shift tab moves to previous tag",
			msg:  tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift},
			assertCmd: func(t *testing.T, cmd tea.Cmd) {
				require.NotNil(t, cmd)
				assert.Equal(t, navigation.TagDeltaMsg{Delta: -1}, cmd())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menu, cleanup := newItemsByTagMenu(t)
			defer cleanup()
			if tt.prepare != nil {
				tt.prepare(&menu)
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
}
