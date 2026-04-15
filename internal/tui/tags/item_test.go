package tags_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rhajizada/donezo/internal/tui/tags"
)

func TestTagItemAccessors(t *testing.T) {
	tests := []struct {
		name            string
		item            tags.Item
		wantTitle       string
		wantFilterValue string
		wantDescription string
	}{
		{
			name:            "plural description",
			item:            tags.NewItem("work", 2),
			wantTitle:       "work",
			wantFilterValue: "work",
			wantDescription: "2 items",
		},
		{
			name:            "singular description",
			item:            tags.NewItem("solo", 1),
			wantTitle:       "solo",
			wantFilterValue: "solo",
			wantDescription: "1 item",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantTitle, tt.item.Title())
			assert.Equal(t, tt.wantFilterValue, tt.item.FilterValue())
			assert.Equal(t, tt.wantDescription, tt.item.Description())
		})
	}

	list := tags.NewList([]tags.Item{tests[0].item, tests[1].item})
	assert.Len(t, list, 2)
}
