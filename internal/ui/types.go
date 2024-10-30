package ui

import (
	"github.com/rhajizada/donezo/pkg/client"
)

// Item wraps repository.Item and implements list.Item interface
type Item struct {
	client.Item
}

// Title returns the title of the item, applying style if completed.
func (i Item) Title() string {
	return i.Item.Title
}

// Description returns the description of the item.
func (i Item) Description() string {
	return i.Item.Description
}

// FilterValue returns the value used for filtering.
func (i Item) FilterValue() string { return i.Item.Title }
