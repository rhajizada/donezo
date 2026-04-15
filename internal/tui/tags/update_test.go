package tags

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/testutil"
)

func TestCopyAndDeleteTag(t *testing.T) {
	tests := []struct {
		name string
		tag  string
	}{
		{name: "copy markdown then delete tag", tag: "work"},
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
			item.Tags = []string{tt.tag}
			_, err = svc.UpdateItem(ctx, item)
			require.NoError(t, err)

			menu := NewModel(ctx, svc)
			msg := menu.ListTags()()
			model, _ := menu.Update(msg)
			menu = model.(MenuModel)
			menu.List.Select(0)

			itemsForTag, err := svc.ListItemsByTag(ctx, tt.tag)
			require.NoError(t, err)
			expected := service.ItemsToMarkdown(tt.tag, *itemsForTag)

			var captured []byte
			prevWrite := writeClipboardText
			writeClipboardText = func(data []byte) { captured = append([]byte{}, data...) }
			defer func() { writeClipboardText = prevWrite }()

			if cmd := menu.Copy(); cmd != nil {
				cmd()
			}
			assert.Equal(t, expected, string(captured))

			delCmd := menu.DeleteTag()
			require.NotNil(t, delCmd)
			msg = delCmd()
			model, _ = menu.Update(msg)
			menu = model.(MenuModel)

			refresh := menu.ListTags()
			msg = refresh()
			model, _ = menu.Update(msg)
			menu = model.(MenuModel)

			assert.Empty(t, menu.List.Items())
			count, err := svc.CountItemsByTag(ctx, tt.tag)
			require.NoError(t, err)
			assert.Zero(t, count)
		})
	}
}
