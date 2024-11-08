package boards

import (
	"github.com/charmbracelet/bubbles/key"
)

// Keymap embeds default list keymap and adds other Binding
type Keymap struct {
	Choose        key.Binding
	CreateBoard   key.Binding
	DeleteBoard   key.Binding
	RenameBoard   key.Binding
	RefreshList   key.Binding
	NextBoard     key.Binding
	PreviousBoard key.Binding
}

func NewKeymap() Keymap {
	return Keymap{
		Choose: key.NewBinding(
			key.WithKeys("return"),
			key.WithHelp("return", "choose board"),
		),
		CreateBoard: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "create a new board"),
		),
		DeleteBoard: key.NewBinding(key.WithKeys("d"),
			key.WithHelp("d", "delete board"),
		),
		RenameBoard: key.NewBinding(key.WithKeys("r"),
			key.WithHelp("r", "rename board"),
		),
		RefreshList: key.NewBinding(key.WithKeys("R"),
			key.WithHelp("R", "refresh list"),
		),
	}
}

func (km Keymap) ShortHelp() []key.Binding {
	bindings := []key.Binding{}
	bindings = append(bindings, km.Choose)
	bindings = append(bindings, km.CreateBoard)
	return bindings
}

func (km Keymap) FullHelp() []key.Binding {
	bindings := []key.Binding{}
	bindings = append(bindings, km.Choose)
	bindings = append(bindings, km.CreateBoard)
	bindings = append(bindings, km.DeleteBoard)
	bindings = append(bindings, km.RenameBoard)
	bindings = append(bindings, km.RefreshList)
	return bindings
}
