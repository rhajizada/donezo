package boards

import (
	"github.com/charmbracelet/bubbles/list"

	"github.com/rhajizada/donezo/internal/service"
)

// Item represents item in the list.
type Item struct {
	Board service.Board
}

func NewList(boards *[]service.Board) []list.Item {
	l := make([]list.Item, len(*boards))
	for i, board := range *boards {
		l[i] = Item{Board: board}
	}
	return l
}

func NewItem(board *service.Board) list.Item {
	return Item{
		Board: *board,
	}
}

func (i Item) Title() string       { return i.Board.Name }
func (i Item) Description() string { return i.Board.CreatedAt.Format("01-02-2006 15:04") }
func (i Item) FilterValue() string { return i.Board.Name }
