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
)

func TestCopySavesItemJSON(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "copy stores serialized item in clipboard"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, cleanup := testutil.NewTestService(t)
			defer cleanup()

			ctx := testutil.MustContext()
			board, err := svc.CreateBoard(ctx, "Inbox")
			require.NoError(t, err)
			item, err := svc.CreateItem(ctx, board, "title", "desc")
			require.NoError(t, err)
			item.Tags = []string{"tag1", "tag2"}
			item.Completed = true
			item, err = svc.UpdateItem(ctx, item)
			require.NoError(t, err)
			items, err := svc.ListItemsByBoard(ctx, board)
			require.NoError(t, err)

			parent := boards.New(ctx, svc)
			parent.List.SetItems(boards.NewList(&[]service.Board{*board}))
			parent.List.Select(0)
			menu := New(ctx, svc, &parent)
			menu.List.SetItems(NewList(items))
			menu.List.Select(0)

			var captured []byte
			prevWrite := writeClipboardText
			writeClipboardText = func(data []byte) { captured = append([]byte{}, data...) }
			defer func() { writeClipboardText = prevWrite }()

			cmd := menu.Copy()
			if cmd != nil {
				cmd()
			}
			assert.NotEmpty(t, captured)

			var saved service.Item
			require.NoError(t, json.Unmarshal(captured, &saved))
			assert.Equal(t, item.Title, saved.Title)
			assert.Equal(t, item.Description, saved.Description)
			assert.True(t, saved.Completed)
			assert.Len(t, saved.Tags, len(item.Tags))
		})
	}
}

func TestPasteCreatesItemOnSelectedBoard(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "paste creates item on selected board"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, cleanup := testutil.NewTestService(t)
			defer cleanup()

			ctx := testutil.MustContext()
			board, err := svc.CreateBoard(ctx, "Inbox")
			require.NoError(t, err)

			parent := boards.New(ctx, svc)
			parent.List.SetItems(boards.NewList(&[]service.Board{*board}))
			parent.List.Select(0)
			menu := New(ctx, svc, &parent)

			clipItem := service.Item{
				Item: repository.Item{
					Title:       "from clipboard",
					Description: "pasted description",
					Completed:   true,
				},
				Tags: []string{"work", "go"},
			}
			data, err := json.Marshal(clipItem)
			require.NoError(t, err)

			prevRead := readClipboardText
			readClipboardText = func() []byte { return data }
			defer func() { readClipboardText = prevRead }()

			cmd := menu.Paste()
			require.NotNil(t, cmd)
			created, ok := cmd().(CreateItemMsg)
			require.True(t, ok)
			require.NoError(t, created.Error)
			require.NotNil(t, created.Item)
			assert.Equal(t, board.ID, created.Item.BoardID)
			assert.Equal(t, clipItem.Title, created.Item.Title)
			assert.Equal(t, clipItem.Description, created.Item.Description)
			assert.True(t, created.Item.Completed)
			assert.Len(t, created.Item.Tags, len(clipItem.Tags))

			items, err := svc.ListItemsByBoard(ctx, board)
			require.NoError(t, err)
			assert.Len(t, *items, 1)
		})
	}
}

func TestItemCreateRenameToggleTagDeleteFlow(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "item create rename toggle update tags and delete flow"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, cleanup := testutil.NewTestService(t)
			defer cleanup()

			ctx := testutil.MustContext()
			board, err := svc.CreateBoard(ctx, "Work")
			require.NoError(t, err)
			parent := boards.New(ctx, svc)
			parent.List.SetItems(boards.NewList(&[]service.Board{*board}))
			parent.List.Select(0)
			menu := New(ctx, svc, &parent)

			menu.InitCreateItem()
			menu.Input.SetValue("Title")
			model, _ := menu.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
			menu = model.(MenuModel)
			assert.Equal(t, CreateItemDescState, menu.Context.State)
			menu.Input.SetValue("Desc")
			model, cmd := menu.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
			menu = model.(MenuModel)
			require.NotNil(t, cmd)
			if msg := cmd(); msg != nil {
				model, _ = menu.Update(msg)
				menu = model.(MenuModel)
			}
			assert.Len(t, menu.List.Items(), 1)

			menu.InitRenameItem()
			menu.Input.SetValue("New Title")
			model, _ = menu.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
			menu = model.(MenuModel)
			assert.Equal(t, RenameItemDescState, menu.Context.State)
			menu.Input.SetValue("New Desc")
			model, cmd = menu.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
			menu = model.(MenuModel)
			require.NotNil(t, cmd)
			if msg := cmd(); msg != nil {
				model, _ = menu.Update(msg)
				menu = model.(MenuModel)
			}
			renamed, ok := menu.selectedItem()
			require.True(t, ok)
			assert.Equal(t, "New Title", renamed.Itm.Title)
			assert.Equal(t, "New Desc", renamed.Itm.Description)

			cmd = menu.ToggleComplete()
			require.NotNil(t, cmd)
			if msg := cmd(); msg != nil {
				model, _ = menu.Update(msg)
				menu = model.(MenuModel)
			}
			toggled, ok := menu.selectedItem()
			require.True(t, ok)
			assert.True(t, toggled.Itm.Completed)

			menu.InitUpdateTags()
			menu.Input.SetValue("one, two")
			model, cmd = menu.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
			menu = model.(MenuModel)
			require.NotNil(t, cmd)
			if msg := cmd(); msg != nil {
				model, _ = menu.Update(msg)
				menu = model.(MenuModel)
			}
			tagged, ok := menu.selectedItem()
			require.True(t, ok)
			assert.Len(t, tagged.Itm.Tags, 2)

			prevWrite := writeClipboardText
			writeClipboardText = func(_ []byte) {}
			defer func() { writeClipboardText = prevWrite }()

			cmd = menu.DeleteItem()
			require.NotNil(t, cmd)
			msg := cmd()
			model, _ = menu.Update(msg)
			menu = model.(MenuModel)
			refresh := menu.ListItems()
			msg = refresh()
			model, _ = menu.Update(msg)
			menu = model.(MenuModel)
			assert.Empty(t, menu.List.Items())

			items, err := svc.ListItemsByBoard(ctx, board)
			require.NoError(t, err)
			assert.Empty(t, *items)
		})
	}
}
