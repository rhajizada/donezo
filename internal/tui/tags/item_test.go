package tags_test

import (
	"testing"

	"github.com/rhajizada/donezo/internal/tui/tags"
)

func TestTagItemAccessors(t *testing.T) {
	item := tags.NewItem("work", 2)
	if item.Title() != "work" {
		t.Fatalf("expected title work, got %q", item.Title())
	}
	if item.FilterValue() != "work" {
		t.Fatalf("expected filter value work, got %q", item.FilterValue())
	}
	if item.Description() != "2 items" {
		t.Fatalf("unexpected plural description %q", item.Description())
	}

	singular := tags.NewItem("solo", 1)
	if singular.Description() != "1 item" {
		t.Fatalf("unexpected singular description %q", singular.Description())
	}

	list := tags.NewList([]tags.Item{item, singular})
	if len(list) != 2 {
		t.Fatalf("expected 2 list items, got %d", len(list))
	}
}
