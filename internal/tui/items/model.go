package items

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rhajizada/donezo/pkg/client"
)

type ItemMenuModel struct {
	Parent  *client.Board
	List    list.Model
	Input   textinput.Model
	Keys    *Keymap
	Context *InputContext
	Client  *client.Client
}

func (m ItemMenuModel) Init() tea.Cmd {
	return m.ListItems()
}

func NewModel(client *client.Client, board *client.Board) ItemMenuModel {
	list := list.New(
		[]list.Item{},
		NewDelegate(),
		0,
		0,
	)
	input := textinput.New()
	keymap := NewKeymap()
	inputContext := NewInputContext()
	list.Title = board.Name
	list.AdditionalShortHelpKeys = keymap.ShortHelp
	list.AdditionalFullHelpKeys = keymap.FullHelp

	return ItemMenuModel{
		Parent:  board,
		List:    list,
		Input:   input,
		Keys:    keymap,
		Context: inputContext,
		Client:  client,
	}
}
