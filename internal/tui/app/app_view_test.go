package app_test

import (
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/rhajizada/donezo/internal/testutil"
	"github.com/rhajizada/donezo/internal/tui/app"
	"github.com/rhajizada/donezo/internal/tui/boards"
)

func TestAppInitAndView(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	ctx := testutil.MustContext()
	m := app.New(ctx, svc)

	initCmd := m.Init()
	if initCmd == nil {
		t.Fatalf("expected init command")
	}
	if msg := initCmd(); msg == nil {
		t.Fatalf("expected init message")
	} else if _, ok := msg.(boards.ListBoardsMsg); !ok {
		t.Fatalf("expected boards.ListBoardsMsg, got %T", msg)
	}

	view := m.View()
	if !view.AltScreen {
		t.Fatalf("expected app view to request alt screen")
	}
}

func TestAppWindowSizeUpdatePath(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	ctx := testutil.MustContext()
	m := app.New(ctx, svc)

	model, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	if _, ok := model.(app.AppModel); !ok {
		t.Fatalf("expected app.AppModel, got %T", model)
	}
}
