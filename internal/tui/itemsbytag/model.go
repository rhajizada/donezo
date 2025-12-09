package itemsbytag

import (
	"context"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/tui/itemlist"
	"github.com/rhajizada/donezo/internal/tui/tags"
)

//nolint:recvcheck // Mixed receivers align with tea.Model usage patterns.
type MenuModel struct {
	ctx     context.Context
	Parent  *tags.MenuModel
	List    itemlist.Model
	Input   textinput.Model
	Keys    *Keymap
	Context *InputContext
	Service *service.Service
}

func (m MenuModel) Init() tea.Cmd {
	return m.ListItems()
}

func New(ctx context.Context, service *service.Service, parent *tags.MenuModel) MenuModel {
	list := itemlist.New(
		[]itemlist.Item{},
		NewDelegate(),
		0,
		0,
	)
	parentItem, ok := parent.List.SelectedItem().(tags.Item)
	if !ok {
		parentItem = tags.Item{}
	}
	input := textinput.New()
	keymap := NewKeymap()
	inputContext := NewInputContext()
	list.Title = parentItem.Tag
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
