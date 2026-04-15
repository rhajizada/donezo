package itemsbytag_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/tui/itemsbytag"
)

func TestItemsByTagItemAccessors(t *testing.T) {
	tests := []struct {
		name       string
		tags       []string
		wantFooter string
	}{
		{name: "tags footer renders list", tags: []string{"work", "go"}, wantFooter: "Tags: work, go"},
		{name: "missing tags footer renders placeholder", tags: nil, wantFooter: "No tags"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := service.Item{Item: service.Item{}.Item, Tags: tt.tags}
			base.Title = "task"
			base.Description = "details"
			base.Completed = true

			item := itemsbytag.NewItem(&base).(itemsbytag.Item)
			assert.Equal(t, "task", item.Title())
			assert.Equal(t, "details", item.Description())
			assert.Equal(t, "task", item.FilterValue())
			assert.True(t, item.HideValue())
			assert.Equal(t, tt.wantFooter, item.Footer())

			list := itemsbytag.NewList(&[]service.Item{base})
			assert.Len(t, list, 1)
		})
	}
}
