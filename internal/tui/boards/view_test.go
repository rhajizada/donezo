package boards_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rhajizada/donezo/internal/testutil"
	"github.com/rhajizada/donezo/internal/tui/boards"
)

func TestBoardsViewForListAndInputStates(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*boards.MenuModel)
		assertView func(*testing.T, string)
	}{
		{
			name: "list state renders content",
			setup: func(menu *boards.MenuModel) {
				menu.List.SetSize(80, 20)
			},
			assertView: func(t *testing.T, view string) {
				assert.NotEmpty(t, strings.TrimSpace(view))
			},
		},
		{
			name: "input state renders typed text",
			setup: func(menu *boards.MenuModel) {
				menu.State = boards.CreateBoardState
				menu.Input.SetValue("new board")
			},
			assertView: func(t *testing.T, view string) {
				assert.Contains(t, view, "new board")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, cleanup := testutil.NewTestService(t)
			defer cleanup()

			ctx := testutil.MustContext()
			menu := boards.New(ctx, svc)
			tt.setup(&menu)
			tt.assertView(t, menu.View().Content)
		})
	}
}
