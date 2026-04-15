package boards

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/testutil"
)

func TestCopyBoardWritesMarkdown(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "copy writes board markdown to clipboard"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, cleanup := testutil.NewTestService(t)
			defer cleanup()

			ctx := testutil.MustContext()
			board, err := svc.CreateBoard(ctx, "Inbox")
			require.NoError(t, err)

			_, err = svc.CreateItem(ctx, board, "task", "desc")
			require.NoError(t, err)
			item2, err := svc.CreateItem(ctx, board, "done", "complete me")
			require.NoError(t, err)
			item2.Completed = true
			_, err = svc.UpdateItem(ctx, item2)
			require.NoError(t, err)

			items, err := svc.ListItemsByBoard(ctx, board)
			require.NoError(t, err)

			menu := New(ctx, svc)
			menu.List.SetItems(NewList(&[]service.Board{*board}))
			menu.List.Select(0)

			expected := service.ItemsToMarkdown(board.Name, *items)

			var captured []byte
			prevWrite := writeClipboardText
			writeClipboardText = func(data []byte) { captured = append([]byte{}, data...) }
			defer func() { writeClipboardText = prevWrite }()

			cmd := menu.Copy()
			if cmd != nil {
				cmd()
			}

			assert.Equal(t, expected, string(captured))
		})
	}
}

func TestBoardCreateRenameDeleteFlow(t *testing.T) {
	tests := []struct {
		name       string
		createName string
		renameTo   string
	}{
		{name: "board create rename and delete flow", createName: "Projects", renameTo: "Renamed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, cleanup := testutil.NewTestService(t)
			defer cleanup()

			ctx := testutil.MustContext()
			menu := New(ctx, svc)

			menu.InitCreateBoard()
			menu.Input.SetValue(tt.createName)
			model, cmd := menu.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
			menu = model.(MenuModel)
			if cmd != nil {
				if msg := cmd(); msg != nil {
					model, _ = menu.Update(msg)
					menu = model.(MenuModel)
				}
			}
			assert.Len(t, menu.List.Items(), 1)
			created := menu.List.SelectedItem().(Item)
			assert.Equal(t, tt.createName, created.Board.Name)

			menu.InitRenameBoard()
			menu.Input.SetValue(tt.renameTo)
			model, cmd = menu.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
			menu = model.(MenuModel)
			if cmd != nil {
				if msg := cmd(); msg != nil {
					model, _ = menu.Update(msg)
					menu = model.(MenuModel)
				}
			}
			renamed := menu.List.SelectedItem().(Item)
			assert.Equal(t, tt.renameTo, renamed.Board.Name)

			cmd = menu.DeleteBoard()
			require.NotNil(t, cmd)
			msg := cmd()
			model, _ = menu.Update(msg)
			menu = model.(MenuModel)

			refresh := menu.ListBoards()
			msg = refresh()
			model, _ = menu.Update(msg)
			menu = model.(MenuModel)
			assert.Empty(t, menu.List.Items())

			boards, err := svc.ListBoards(ctx)
			require.NoError(t, err)
			assert.Empty(t, *boards)
		})
	}
}
