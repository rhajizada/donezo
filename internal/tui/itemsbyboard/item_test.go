package itemsbyboard_test

import (
	"testing"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/tui/itemsbyboard"
)

func TestItemsByBoardItemAccessors(t *testing.T) {
	base := service.Item{Item: service.Item{}.Item, Tags: []string{"work"}}
	base.Title = "task"
	base.Description = "line 1\nline 2"

	item := itemsbyboard.NewItem(&base).(itemsbyboard.Item)
	if item.Title() != "task" {
		t.Fatalf("expected title task, got %q", item.Title())
	}
	if item.Description() == "" {
		t.Fatalf("expected non-empty description")
	}
	if item.FilterValue() != "task" {
		t.Fatalf("expected filter value task, got %q", item.FilterValue())
	}
	if item.HideValue() {
		t.Fatalf("expected incomplete item to remain visible")
	}
	if item.Footer() != "Tags: work" {
		t.Fatalf("unexpected footer %q", item.Footer())
	}

	base.Completed = true
	item = itemsbyboard.NewItem(&base).(itemsbyboard.Item)
	if !item.HideValue() {
		t.Fatalf("expected completed item to be hidden")
	}

	list := itemsbyboard.NewList(&[]service.Item{base})
	if len(list) != 1 {
		t.Fatalf("expected 1 list item, got %d", len(list))
	}
}
