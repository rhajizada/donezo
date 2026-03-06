package itemsbytag_test

import (
	"testing"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/tui/itemsbytag"
)

func TestItemsByTagItemAccessors(t *testing.T) {
	base := service.Item{Item: service.Item{}.Item, Tags: []string{"work", "go"}}
	base.Title = "task"
	base.Description = "details"
	base.Completed = true

	item := itemsbytag.NewItem(&base).(itemsbytag.Item)
	if item.Title() != "task" {
		t.Fatalf("expected title task, got %q", item.Title())
	}
	if item.Description() != "details" {
		t.Fatalf("expected description details, got %q", item.Description())
	}
	if item.FilterValue() != "task" {
		t.Fatalf("expected filter value task, got %q", item.FilterValue())
	}
	if !item.HideValue() {
		t.Fatalf("expected completed item to be hidden")
	}
	if item.Footer() != "Tags: work, go" {
		t.Fatalf("unexpected footer %q", item.Footer())
	}

	base.Tags = nil
	item = itemsbytag.NewItem(&base).(itemsbytag.Item)
	if item.Footer() != "No tags" {
		t.Fatalf("unexpected footer without tags: %q", item.Footer())
	}

	list := itemsbytag.NewList(&[]service.Item{base})
	if len(list) != 1 {
		t.Fatalf("expected 1 list item, got %d", len(list))
	}
}
