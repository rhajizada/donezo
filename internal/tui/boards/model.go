package boards

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rhajizada/donezo/pkg/client"
)

type MenuModel struct {
	List   list.Model
	Input  textinput.Model
	Keys   *Keymap
	State  InputState
	Client *client.Client
}

func (m MenuModel) Init() tea.Cmd {
	return m.ListBoards()
}

func NewModel(client *client.Client) MenuModel {
	list := list.New(
		[]list.Item{},
		list.NewDefaultDelegate(),
		0,
		0,
	)
	input := textinput.New()
	keymap := NewKeymap()
	list.Title = "donezo"
	list.AdditionalShortHelpKeys = keymap.ShortHelp
	list.AdditionalFullHelpKeys = keymap.FullHelp
	return MenuModel{
		List:   list,
		Input:  input,
		State:  DefaultState,
		Keys:   &keymap,
		Client: client,
	}
}
