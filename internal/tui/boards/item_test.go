package boards_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/testutil"
	"github.com/rhajizada/donezo/internal/tui/boards"
)

func TestBoardItemAccessors(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "board item exposes accessors"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, cleanup := testutil.NewTestService(t)
			defer cleanup()

			ctx := testutil.MustContext()
			board, err := svc.CreateBoard(ctx, "Inbox")
			require.NoError(t, err)

			item := boards.NewItem(board).(boards.Item)
			assert.Equal(t, "Inbox", item.Title())
			assert.Equal(t, "Inbox", item.FilterValue())
			assert.NotEmpty(t, item.Description())

			list := boards.NewList(&[]service.Board{*board})
			assert.Len(t, list, 1)
		})
	}
}
