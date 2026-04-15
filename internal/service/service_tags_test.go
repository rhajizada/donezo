package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/testutil"
)

func TestListItemsByTagListTagsAndDeleteTag(t *testing.T) {
	tests := []struct {
		name                string
		itemATags           []string
		itemBTags           []string
		queryTag            string
		wantItemsByTagCount int
		wantTagCount        int
	}{
		{
			name:                "list by tag list tags and delete tag",
			itemATags:           []string{"work", "go"},
			itemBTags:           []string{"work"},
			queryTag:            "work",
			wantItemsByTagCount: 2,
			wantTagCount:        2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, cleanup := testutil.NewTestService(t)
			defer cleanup()

			ctx := testutil.MustContext()
			board := mustCreateBoard(ctx, t, svc, "Coverage")
			itemA := mustCreateItem(ctx, t, svc, board, "a", "first")
			itemB := mustCreateItem(ctx, t, svc, board, "b", "second")

			itemA.Tags = tt.itemATags
			_ = mustUpdateItem(ctx, t, svc, itemA)
			itemB.Tags = tt.itemBTags
			_ = mustUpdateItem(ctx, t, svc, itemB)

			itemsByTag, err := svc.ListItemsByTag(ctx, tt.queryTag)
			require.NoError(t, err)
			assert.Len(t, *itemsByTag, tt.wantItemsByTagCount)

			tags, err := svc.ListTags(ctx)
			require.NoError(t, err)
			assert.Len(t, tags, tt.wantTagCount)

			require.NoError(t, svc.DeleteTag(ctx, tt.queryTag))
			count, err := svc.CountItemsByTag(ctx, tt.queryTag)
			require.NoError(t, err)
			assert.Zero(t, count)
		})
	}
}

func TestItemsToMarkdownFormatting(t *testing.T) {
	items := []service.Item{
		{Item: service.Item{}.Item, Tags: []string{"work"}},
		{Item: service.Item{}.Item, Tags: []string{"go"}},
	}
	items[0].Title = "open"
	items[0].Description = "desc one"
	items[1].Title = "done"
	items[1].Description = "desc two"
	items[1].Completed = true

	tests := []struct {
		name   string
		needle string
	}{
		{name: "renders markdown header", needle: "### Today"},
		{name: "renders unchecked line", needle: "- [ ] **open**"},
		{name: "renders checked line", needle: "- [X] **done**"},
		{name: "renders description line", needle: "\t- desc two"},
	}

	md := service.ItemsToMarkdown("Today", items)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Contains(t, md, tt.needle)
		})
	}
}

func TestUpdateItemRejectsEmptyTags(t *testing.T) {
	tests := []struct {
		name   string
		tags   []string
		needle string
	}{
		{name: "rejects empty tag", tags: []string{"ok", ""}, needle: "tag must not be empty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, cleanup := testutil.NewTestService(t)
			defer cleanup()

			ctx := testutil.MustContext()
			board := mustCreateBoard(ctx, t, svc, "Validation")
			item := mustCreateItem(ctx, t, svc, board, "task", "desc")
			item.Tags = tt.tags

			updated, err := svc.UpdateItem(ctx, item)
			assert.Nil(t, updated)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.needle)
		})
	}
}
