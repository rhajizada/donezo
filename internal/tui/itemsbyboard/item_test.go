package itemsbyboard_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/tui/itemsbyboard"
)

func TestItemsByBoardItemAccessors(t *testing.T) {
	tests := []struct {
		name       string
		completed  bool
		wantHidden bool
		wantFooter string
	}{
		{name: "incomplete item stays visible", completed: false, wantHidden: false, wantFooter: "Tags: work"},
		{name: "completed item is hidden", completed: true, wantHidden: true, wantFooter: "Tags: work"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := service.Item{Item: service.Item{}.Item, Tags: []string{"work"}}
			base.Title = "task"
			base.Description = "line 1\nline 2"
			base.Completed = tt.completed

			item := itemsbyboard.NewItem(&base).(itemsbyboard.Item)
			assert.Equal(t, "task", item.Title())
			assert.NotEmpty(t, item.Description())
			assert.Equal(t, "task", item.FilterValue())
			assert.Equal(t, tt.wantHidden, item.HideValue())
			assert.Equal(t, tt.wantFooter, item.Footer())

			list := itemsbyboard.NewList(&[]service.Item{base})
			assert.Len(t, list, 1)
		})
	}
}
