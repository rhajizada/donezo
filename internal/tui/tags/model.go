package tags

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rhajizada/donezo/internal/service"
)

type MenuModel struct {
	ctx    context.Context
	List   list.Model
	Keys   *Keymap
	Client *service.Service
}

func (m MenuModel) Init() tea.Cmd {
	return m.ListTags()
}

func NewModel(ctx context.Context, client *service.Service) MenuModel {
	list := list.New(
		[]list.Item{},
		list.NewDefaultDelegate(),
		0,
		0,
	)
	keymap := NewKeymap()
	list.Title = "donezo | Tags"
	list.AdditionalShortHelpKeys = keymap.ShortHelp
	list.AdditionalFullHelpKeys = keymap.FullHelp
	return MenuModel{
		ctx:    ctx,
		List:   list,
		Keys:   &keymap,
		Client: client,
	}
}
