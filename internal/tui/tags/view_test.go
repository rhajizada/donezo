package tags_test

import (
	"strings"
	"testing"

	"github.com/rhajizada/donezo/internal/testutil"
	"github.com/rhajizada/donezo/internal/tui/tags"
)

func TestTagsViewRendersListTitle(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	ctx := testutil.MustContext()
	menu := tags.NewModel(ctx, svc)
	menu.List.SetSize(80, 20)

	view := menu.View().Content
	if strings.TrimSpace(view) == "" {
		t.Fatalf("expected non-empty tags view")
	}
}
