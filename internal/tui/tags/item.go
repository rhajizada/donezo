package tags

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
)

// Item represents item in the list.
type Item struct {
	Tag   string
	Count int64
}

func NewList(tags []Item) []list.Item {
	l := make([]list.Item, len(tags))
	for i, tag := range tags {
		l[i] = tag
	}
	return l
}

func NewItem(tag string, count int64) Item {
	return Item{
		Tag:   tag,
		Count: count,
	}
}

func (i Item) Title() string { return i.Tag }
func (i Item) Description() string {
	var suffix string
	if i.Count != 1 {
		suffix = "s"
	}
	return fmt.Sprintf("%d item%s", i.Count, suffix)
}
func (i Item) FilterValue() string { return i.Tag }
