package itemsbytag_test

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/testutil"
	"github.com/rhajizada/donezo/internal/tui/itemsbytag"
	"github.com/rhajizada/donezo/internal/tui/tags"
)

func newItemsByTagMenu(t *testing.T) (itemsbytag.MenuModel, func()) {
	t.Helper()
	svc, cleanup := testutil.NewTestService(t)
	ctx := testutil.MustContext()

	board, err := svc.CreateBoard(ctx, "Inbox")
	require.NoError(t, err)
	item, err := svc.CreateItem(ctx, board, "task", "desc")
	require.NoError(t, err)
	item.Tags = []string{"work"}
	item, err = svc.UpdateItem(ctx, item)
	require.NoError(t, err)

	tagCount, err := svc.CountItemsByTag(ctx, "work")
	require.NoError(t, err)
	parent := tags.NewModel(ctx, svc)
	parent.List.SetItems(tags.NewList([]tags.Item{tags.NewItem("work", tagCount)}))
	parent.List.Select(0)

	menu := itemsbytag.New(ctx, svc, &parent)
	menu.List.SetItems(itemsbytag.NewList(&[]service.Item{*item}))
	menu.List.Select(0)
	return menu, cleanup
}

func TestListItemsWithoutSelectedTagReturnsErrorMsg(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "missing selected tag returns error message"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menu, cleanup := newItemsByTagMenu(t)
			defer cleanup()

			menu.Parent.List.SetItems(tags.NewList(nil))
			_, ok := menu.ListItems()().(itemsbytag.ErrorMsg)
			assert.True(t, ok)
		})
	}
}

func TestRenameUpdateTagsAndToggleCommands(t *testing.T) {
	tests := []struct {
		name string
		run  func(*testing.T, itemsbytag.MenuModel)
	}{
		{
			name: "rename item succeeds",
			run: func(t *testing.T, menu itemsbytag.MenuModel) {
				menu.Context.Title = "renamed"
				menu.Context.Desc = "renamed desc"
				renameMsg, ok := menu.RenameItem()().(itemsbytag.RenameItemMsg)
				require.True(t, ok)
				assert.NoError(t, renameMsg.Error)
			},
		},
		{
			name: "update tags succeeds",
			run: func(t *testing.T, menu itemsbytag.MenuModel) {
				menu.Context.Title = "one, two"
				updateTagsMsg, ok := menu.UpdateTags()().(itemsbytag.UpdateTagsMsg)
				require.True(t, ok)
				require.NoError(t, updateTagsMsg.Error)
				assert.Len(t, updateTagsMsg.Item.Tags, 2)
			},
		},
		{
			name: "toggle complete succeeds",
			run: func(t *testing.T, menu itemsbytag.MenuModel) {
				toggleMsg, ok := menu.ToggleComplete()().(itemsbytag.ToggleItemMsg)
				require.True(t, ok)
				require.NoError(t, toggleMsg.Error)
				assert.True(t, toggleMsg.Item.Completed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menu, cleanup := newItemsByTagMenu(t)
			defer cleanup()
			tt.run(t, menu)
		})
	}
}

func TestUpdateTagsValidationError(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "invalid tag list returns validation error", input: "ok, "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menu, cleanup := newItemsByTagMenu(t)
			defer cleanup()

			menu.Context.Title = tt.input
			updateMsg, ok := menu.UpdateTags()().(itemsbytag.UpdateTagsMsg)
			require.True(t, ok)
			assert.Error(t, updateMsg.Error)
		})
	}
}

func TestHandleInputStateTransitions(t *testing.T) {
	tests := []struct {
		name      string
		state     itemsbytag.InputState
		value     string
		msg       tea.KeyPressMsg
		wantState itemsbytag.InputState
		checkCmds bool
		wantCmds  bool
	}{
		{
			name:      "rename title advances to description",
			state:     itemsbytag.RenameItemNameState,
			value:     "new title",
			msg:       tea.KeyPressMsg{Code: tea.KeyEnter},
			wantState: itemsbytag.RenameItemDescState,
		},
		{
			name:      "rename description completes flow",
			state:     itemsbytag.RenameItemDescState,
			value:     "new description",
			msg:       tea.KeyPressMsg{Code: tea.KeyEnter},
			wantState: itemsbytag.DefaultState,
			checkCmds: true,
			wantCmds:  true,
		},
		{
			name:      "update tags completes flow",
			state:     itemsbytag.UpdateTagsState,
			value:     "a, b",
			msg:       tea.KeyPressMsg{Code: tea.KeyEnter},
			wantState: itemsbytag.DefaultState,
			checkCmds: true,
			wantCmds:  true,
		},
		{
			name:      "escape cancels input flow",
			state:     itemsbytag.RenameItemNameState,
			msg:       tea.KeyPressMsg{Code: tea.KeyEsc},
			wantState: itemsbytag.DefaultState,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menu, cleanup := newItemsByTagMenu(t)
			defer cleanup()

			menu.Context.State = tt.state
			menu.Input.SetValue(tt.value)
			_, cmds := menu.HandleInputState(tt.msg)
			assert.Equal(t, tt.wantState, menu.Context.State)
			if tt.checkCmds {
				if tt.wantCmds {
					assert.NotEmpty(t, cmds)
				} else {
					assert.Empty(t, cmds)
				}
			}
		})
	}
}

func TestUpdateWithListItemsMsgReplacesList(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "list items message replaces current list"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menu, cleanup := newItemsByTagMenu(t)
			defer cleanup()

			newItems := &[]service.Item{{Item: service.Item{}.Item, Tags: []string{"x"}}}
			(*newItems)[0].Title = "replacement"
			(*newItems)[0].Description = "desc"

			model, _ := menu.Update(itemsbytag.ListItemsMsg{Items: newItems})
			updated := model.(itemsbytag.MenuModel)
			assert.Len(t, updated.List.Items(), 1)
		})
	}
}
