package itemsbytag_test

import (
	"strings"
	"testing"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/testutil"
	"github.com/rhajizada/donezo/internal/tui/itemsbytag"
	"github.com/rhajizada/donezo/internal/tui/tags"
)

func TestItemsByTagViewForListAndInputStates(t *testing.T) {
	svc, cleanup := testutil.NewTestService(t)
	defer cleanup()

	ctx := testutil.MustContext()
	board, err := svc.CreateBoard(ctx, "Inbox")
	if err != nil {
		t.Fatalf("CreateBoard: %v", err)
	}
	item, err := svc.CreateItem(ctx, board, "task", "desc")
	if err != nil {
		t.Fatalf("CreateItem: %v", err)
	}
	item.Tags = []string{"work"}
	_, err = svc.UpdateItem(ctx, item)
	if err != nil {
		t.Fatalf("UpdateItem: %v", err)
	}

	tagCount, err := svc.CountItemsByTag(ctx, "work")
	if err != nil {
		t.Fatalf("CountItemsByTag: %v", err)
	}
	parent := tags.NewModel(ctx, svc)
	parent.List.SetItems(tags.NewList([]tags.Item{tags.NewItem("work", tagCount)}))
	parent.List.Select(0)

	menu := itemsbytag.New(ctx, svc, &parent)
	menu.List.SetItems(itemsbytag.NewList(&[]service.Item{*item}))
	menu.List.Select(0)
	menu.List.SetSize(80, 20)

	listView := menu.View().Content
	if strings.TrimSpace(listView) == "" {
		t.Fatalf("expected non-empty items-by-tag list view")
	}

	menu.Context.State = itemsbytag.UpdateTagsState
	menu.Input.SetValue("a, b")
	inputView := menu.View().Content
	if !strings.Contains(inputView, "a, b") {
		t.Fatalf("expected input text in input view")
	}
}
