package service_test

import (
	"strings"
	"testing"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/testutil"
)

func TestListItemsByTagListTagsAndDeleteTag(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	ctx := testutil.MustContext()
	board := mustCreateBoard(ctx, t, svc, "Coverage")
	itemA := mustCreateItem(ctx, t, svc, board, "a", "first")
	itemB := mustCreateItem(ctx, t, svc, board, "b", "second")

	itemA.Tags = []string{"work", "go"}
	_ = mustUpdateItem(ctx, t, svc, itemA)
	itemB.Tags = []string{"work"}
	_ = mustUpdateItem(ctx, t, svc, itemB)

	itemsByTag, err := svc.ListItemsByTag(ctx, "work")
	requireNoError(t, err, "ListItemsByTag")
	if len(*itemsByTag) != 2 {
		t.Fatalf("expected 2 items for work tag, got %d", len(*itemsByTag))
	}

	tags, err := svc.ListTags(ctx)
	requireNoError(t, err, "ListTags")
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}

	requireNoError(t, svc.DeleteTag(ctx, "work"), "DeleteTag")
	count, err := svc.CountItemsByTag(ctx, "work")
	requireNoError(t, err, "CountItemsByTag after delete")
	if count != 0 {
		t.Fatalf("expected tag count 0 after delete, got %d", count)
	}
}

func TestItemsToMarkdownFormatting(t *testing.T) {
	items := []service.Item{
		{
			Item: service.Item{}.Item,
			Tags: []string{"work"},
		},
		{
			Item: service.Item{}.Item,
			Tags: []string{"go"},
		},
	}
	items[0].Title = "open"
	items[0].Description = "desc one"
	items[1].Title = "done"
	items[1].Description = "desc two"
	items[1].Completed = true

	md := service.ItemsToMarkdown("Today", items)
	if !strings.Contains(md, "### Today") {
		t.Fatalf("expected markdown header, got %q", md)
	}
	if !strings.Contains(md, "- [ ] **open**") {
		t.Fatalf("expected unchecked markdown line, got %q", md)
	}
	if !strings.Contains(md, "- [X] **done**") {
		t.Fatalf("expected checked markdown line, got %q", md)
	}
	if !strings.Contains(md, "\t- desc two") {
		t.Fatalf("expected description line, got %q", md)
	}
}

func TestUpdateItemRejectsEmptyTags(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	ctx := testutil.MustContext()
	board := mustCreateBoard(ctx, t, svc, "Validation")
	item := mustCreateItem(ctx, t, svc, board, "task", "desc")
	item.Tags = []string{"ok", ""}

	updated, err := svc.UpdateItem(ctx, item)
	if err == nil {
		t.Fatalf("expected validation error, got updated item %+v", updated)
	}
	if !strings.Contains(err.Error(), "tag must not be empty") {
		t.Fatalf("unexpected error %v", err)
	}
}
