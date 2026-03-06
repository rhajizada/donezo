package tags

import (
	"context"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"

	"github.com/rhajizada/donezo/internal/service"
)

//nolint:recvcheck // Mixed receivers align with tea.Model usage patterns.
type MenuModel struct {
	ctx    context.Context
	List   list.Model
	Keys   *Keymap
	Client *service.Service
}

// NewModel constructs a new tag list menu.
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

func (m MenuModel) Init() tea.Cmd {
	return m.ListTags()
}
