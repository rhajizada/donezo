package tags_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rhajizada/donezo/internal/testutil"
	"github.com/rhajizada/donezo/internal/tui/tags"
)

func TestTagsViewRendersListTitle(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "renders non-empty tags view"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, cleanup := testutil.NewTestService(t)
			defer cleanup()

			ctx := testutil.MustContext()
			menu := tags.NewModel(ctx, svc)
			menu.List.SetSize(80, 20)

			view := menu.View().Content
			assert.NotEmpty(t, strings.TrimSpace(view))
		})
	}
}
