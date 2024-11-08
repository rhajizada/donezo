package items

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/rhajizada/donezo/pkg/client"
)

// Item represents item in the list
type Item struct {
	Itm client.Item
}

func NewList(items *[]client.Item) []list.Item {
	l := make([]list.Item, len(*items))
	for i, item := range *items {
		l[i] = Item{Itm: item}
	}
	return l
}

func NewItem(item *client.Item) list.Item {
	return Item{
		Itm: *item,
	}
}

func (i Item) Title() string       { return i.Itm.Title }
func (i Item) Description() string { return i.Itm.Description }
func (i Item) FilterValue() string { return i.Itm.Title }
