package boards_test

import (
	"strings"
	"testing"

	"github.com/rhajizada/donezo/internal/testutil"
	"github.com/rhajizada/donezo/internal/tui/boards"
)

func TestBoardsViewForListAndInputStates(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	ctx := testutil.MustContext()
	menu := boards.New(ctx, svc)
	menu.List.SetSize(80, 20)

	listView := menu.View().Content
	if strings.TrimSpace(listView) == "" {
		t.Fatalf("expected non-empty boards list view")
	}

	menu.State = boards.CreateBoardState
	menu.Input.SetValue("new board")
	inputView := menu.View().Content
	if !strings.Contains(inputView, "new board") {
		t.Fatalf("expected input text in input view")
	}
}
