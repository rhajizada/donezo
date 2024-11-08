package boards

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/rhajizada/donezo/pkg/client"
)

// Item represents item in the list
type Item struct {
	Board client.Board
}

func NewList(boards *[]client.Board) []list.Item {
	l := make([]list.Item, len(*boards))
	for i, board := range *boards {
		l[i] = Item{Board: board}
	}
	return l
}

func NewItem(board *client.Board) list.Item {
	return Item{
		Board: *board,
	}
}

func (i Item) Title() string       { return i.Board.Name }
func (i Item) Description() string { return i.Board.CreatedAt.Format("01-02-2006 15:04") }
func (i Item) FilterValue() string { return i.Board.Name }
