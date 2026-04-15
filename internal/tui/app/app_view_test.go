package app_test

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rhajizada/donezo/internal/testutil"
	"github.com/rhajizada/donezo/internal/tui/app"
	"github.com/rhajizada/donezo/internal/tui/boards"
)

func TestAppInitAndView(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "init returns boards message and view uses alt screen"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, cleanup := testutil.NewTestService(t)
			defer cleanup()

			ctx := testutil.MustContext()
			m := app.New(ctx, svc)

			initCmd := m.Init()
			require.NotNil(t, initCmd)
			msg := initCmd()
			require.NotNil(t, msg)
			_, ok := msg.(boards.ListBoardsMsg)
			assert.True(t, ok)

			view := m.View()
			assert.True(t, view.AltScreen)
		})
	}
}

func TestAppWindowSizeUpdatePath(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{name: "window size update returns app model", width: 120, height: 40},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, cleanup := testutil.NewTestService(t)
			defer cleanup()

			ctx := testutil.MustContext()
			m := app.New(ctx, svc)

			model, _ := m.Update(tea.WindowSizeMsg{Width: tt.width, Height: tt.height})
			_, ok := model.(app.AppModel)
			assert.True(t, ok)
		})
	}
}
