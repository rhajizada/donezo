package items

import "github.com/charmbracelet/bubbles/key"

type Keymap struct {
	Back           key.Binding
	CreateItem     key.Binding
	DeleteItem     key.Binding
	RenameItem     key.Binding
	RefreshList    key.Binding
	ToggleComplete key.Binding
	NextBoard      key.Binding
	PreviousBoard  key.Binding
}

func NewKeymap() *Keymap {
	return &Keymap{
		Back: key.NewBinding(
			key.WithKeys("b"),
			key.WithHelp("b", "back"),
		),
		CreateItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "create item"),
		),
		DeleteItem: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete item"),
		),
		RenameItem: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "rename item"),
		),
		RefreshList: key.NewBinding(
			key.WithKeys("R"),
			key.WithHelp("R", "refresh board"),
		),
		ToggleComplete: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle complete"),
		),
		NextBoard: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next board"),
		),
		PreviousBoard: key.NewBinding(
			key.WithKeys("shift + tab"),
			key.WithHelp("shift+tab", "previous board"),
		),
	}
}

func (km Keymap) ShortHelp() []key.Binding {
	bindings := []key.Binding{}
	bindings = append(bindings, km.Back)
	bindings = append(bindings, km.CreateItem)
	return bindings
}

func (km Keymap) FullHelp() []key.Binding {
	bindings := []key.Binding{}
	bindings = append(bindings, km.Back)
	bindings = append(bindings, km.CreateItem)
	bindings = append(bindings, km.DeleteItem)
	bindings = append(bindings, km.RenameItem)
	bindings = append(bindings, km.RefreshList)
	bindings = append(bindings, km.ToggleComplete)
	bindings = append(bindings, km.NextBoard)
	bindings = append(bindings, km.PreviousBoard)
	return bindings
}
