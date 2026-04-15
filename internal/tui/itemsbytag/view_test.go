package itemsbytag_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/testutil"
	"github.com/rhajizada/donezo/internal/tui/itemsbytag"
	"github.com/rhajizada/donezo/internal/tui/tags"
)

func TestItemsByTagViewForListAndInputStates(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*itemsbytag.MenuModel)
		assertView func(*testing.T, string)
	}{
		{
			name: "list state renders content",
			setup: func(menu *itemsbytag.MenuModel) {
				menu.List.SetSize(80, 20)
			},
			assertView: func(t *testing.T, view string) {
				assert.NotEmpty(t, strings.TrimSpace(view))
			},
		},
		{
			name: "input state renders typed text",
			setup: func(menu *itemsbytag.MenuModel) {
				menu.Context.State = itemsbytag.UpdateTagsState
				menu.Input.SetValue("a, b")
			},
			assertView: func(t *testing.T, view string) {
				assert.Contains(t, view, "a, b")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, cleanup := testutil.NewTestService(t)
			defer cleanup()

			ctx := testutil.MustContext()
			board, err := svc.CreateBoard(ctx, "Inbox")
			require.NoError(t, err)
			item, err := svc.CreateItem(ctx, board, "task", "desc")
			require.NoError(t, err)
			item.Tags = []string{"work"}
			_, err = svc.UpdateItem(ctx, item)
			require.NoError(t, err)

			tagCount, err := svc.CountItemsByTag(ctx, "work")
			require.NoError(t, err)
			parent := tags.NewModel(ctx, svc)
			parent.List.SetItems(tags.NewList([]tags.Item{tags.NewItem("work", tagCount)}))
			parent.List.Select(0)

			menu := itemsbytag.New(ctx, svc, &parent)
			menu.List.SetItems(itemsbytag.NewList(&[]service.Item{*item}))
			menu.List.Select(0)
			tt.setup(&menu)
			tt.assertView(t, menu.View().Content)
		})
	}
}
