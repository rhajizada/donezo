package itemsbyboard_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/testutil"
	"github.com/rhajizada/donezo/internal/tui/boards"
	"github.com/rhajizada/donezo/internal/tui/itemsbyboard"
)

func TestItemsByBoardViewForListAndInputStates(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*itemsbyboard.MenuModel)
		assertView func(*testing.T, string)
	}{
		{
			name: "list state renders content",
			setup: func(menu *itemsbyboard.MenuModel) {
				menu.List.SetSize(80, 20)
				menu.List.Title = "Inbox"
			},
			assertView: func(t *testing.T, view string) {
				assert.NotEmpty(t, strings.TrimSpace(view))
			},
		},
		{
			name: "input state renders typed text",
			setup: func(menu *itemsbyboard.MenuModel) {
				menu.Context.State = itemsbyboard.CreateItemNameState
				menu.Input.SetValue("new item")
			},
			assertView: func(t *testing.T, view string) {
				assert.Contains(t, view, "new item")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, cleanup := testutil.NewTestService(t)
			defer cleanup()

			ctx := testutil.MustContext()
			boardMenu := boards.New(ctx, svc)
			board, err := svc.CreateBoard(ctx, "Inbox")
			require.NoError(t, err)
			boardMenu.List.SetItems(boards.NewList(&[]service.Board{*board}))
			boardMenu.List.Select(0)

			menu := itemsbyboard.New(ctx, svc, &boardMenu)
			tt.setup(&menu)
			tt.assertView(t, menu.View().Content)
		})
	}
}
