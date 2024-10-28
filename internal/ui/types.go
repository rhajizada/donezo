package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/rhajizada/donezo/internal/repository"
)

// Define a style for completed items
var completedItemStyle = lipgloss.NewStyle().Strikethrough(true)

// Item wraps repository.Item and implements list.Item interface
type Item struct {
	repository.Item
}

// Title returns the title of the item, applying style if completed.
func (i Item) Title() string {
	// TODO:
	// Fix it so that strike thtrough is added to the existing style
	if i.Completed {
		return completedItemStyle.Render(i.Item.Title)
	}
	return i.Item.Title
}

// Description returns the description of the item.
func (i Item) Description() string {
	if i.Completed {
		return completedItemStyle.Render(i.Item.Description)
	}
	return i.Item.Description
}

// FilterValue returns the value used for filtering.
func (i Item) FilterValue() string { return i.Item.Title }
