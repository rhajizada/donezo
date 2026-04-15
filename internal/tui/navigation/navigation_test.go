package navigation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViewOrdering(t *testing.T) {
	tests := []struct {
		name string
		view View
		want int
	}{
		{name: "boards view", view: ViewBoards, want: 0},
		{name: "tags view", view: ViewTags, want: 1},
		{name: "items by board view", view: ViewItemsByBoard, want: 2},
		{name: "items by tag view", view: ViewItemsByTag, want: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, int(tt.view))
		})
	}
}
