package itemsbyboard

import (
	"context"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/tui/boards"
	"github.com/rhajizada/donezo/internal/tui/itemlist"
)

type MenuModel struct {
	ctx     context.Context
	Parent  *boards.MenuModel
	List    itemlist.Model
	Input   textinput.Model
	Keys    *Keymap
	Context *InputContext
	Service *service.Service
}

func (m MenuModel) Init() tea.Cmd {
	return m.ListItems()
}

func New(ctx context.Context, service *service.Service, parent *boards.MenuModel) MenuModel {
	list := itemlist.New(
		[]itemlist.Item{},
		NewDelegate(),
		0,
		0,
	)
	parentItem := parent.List.SelectedItem().(boards.Item)
	input := textinput.New()
	keymap := NewKeymap()
	inputContext := NewInputContext()
	list.Title = parentItem.Board.Name
	list.AdditionalShortHelpKeys = keymap.ShortHelp
	list.AdditionalFullHelpKeys = keymap.FullHelp

	return MenuModel{
		ctx:     ctx,
		Parent:  parent,
		List:    list,
		Input:   input,
		Keys:    keymap,
		Context: inputContext,
		Service: service,
	}
}
