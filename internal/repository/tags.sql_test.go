package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rhajizada/donezo/internal/repository"
)

func TestTagQueries(t *testing.T) {
	tests := []struct {
		name string
		run  func(*testing.T, *repository.Queries)
	}{
		{
			name: "add list count remove and delete tags",
			run: func(t *testing.T, q *repository.Queries) {
				ctx := context.Background()
				board := mustCreateBoard(t, q, "Inbox")
				itemA := mustCreateItem(t, q, board.ID, "a", "first")
				itemB := mustCreateItem(t, q, board.ID, "b", "second")

				for _, tag := range []repository.AddTagToItemByIDParams{
					{ItemID: itemA.ID, Tag: "work"},
					{ItemID: itemA.ID, Tag: "go"},
					{ItemID: itemB.ID, Tag: "work"},
				} {
					require.NoError(t, q.AddTagToItemByID(ctx, tag))
				}

				count, err := q.CountItemsByTag(ctx, "work")
				require.NoError(t, err)
				assert.EqualValues(t, 2, count)

				tags, err := q.ListTags(ctx)
				require.NoError(t, err)
				assert.Equal(t, []string{"go", "work"}, tags)

				itemTags, err := q.ListTagsByItemID(ctx, itemA.ID)
				require.NoError(t, err)
				assert.Equal(t, []string{"go", "work"}, itemTags)

				itemsByTag, err := q.ListItemsByTag(ctx, "work")
				require.NoError(t, err)
				require.Len(t, itemsByTag, 2)
				assert.Contains(t, tagsJSON(t, itemsByTag[0].Tags), "work")

				require.NoError(t, q.RemoveTagFromItemByID(ctx, repository.RemoveTagFromItemByIDParams{
					ItemID: itemA.ID,
					Tag:    "work",
				}))

				count, err = q.CountItemsByTag(ctx, "work")
				require.NoError(t, err)
				assert.EqualValues(t, 1, count)

				require.NoError(t, q.DeleteTag(ctx, "work"))
				count, err = q.CountItemsByTag(ctx, "work")
				require.NoError(t, err)
				assert.Zero(t, count)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, q := newTestQueries(t)
			tt.run(t, q)
		})
	}
}
