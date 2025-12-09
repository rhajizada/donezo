package itemsbyboard

import "github.com/charmbracelet/bubbles/key"

type Keymap struct {
	Back           key.Binding
	CreateItem     key.Binding
	DeleteItem     key.Binding
	RenameItem     key.Binding
	UpdateTags     key.Binding
	RefreshList    key.Binding
	ToggleComplete key.Binding
	Cut            key.Binding
	Copy           key.Binding
	Paste          key.Binding
	NextBoard      key.Binding
	PreviousBoard  key.Binding
}

func NewKeymap() *Keymap {
	return &Keymap{
		Back: key.NewBinding(
			key.WithKeys("backspace"),
			key.WithHelp("backspace", "back"),
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
		UpdateTags: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "update tags"),
		),
		RefreshList: key.NewBinding(
			key.WithKeys("R"),
			key.WithHelp("R", "refresh board"),
		),
		ToggleComplete: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle complete"),
		),
		Copy: key.NewBinding(
			key.WithKeys("y"),
			key.WithKeys("y", "copy"),
		),
		Paste: key.NewBinding(
			key.WithKeys("p"),
			key.WithKeys("p", "paste"),
		),
		NextBoard: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next board"),
		),
		PreviousBoard: key.NewBinding(
			key.WithKeys("shift+tab"),
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
	bindings = append(bindings, km.UpdateTags)
	bindings = append(bindings, km.RefreshList)
	bindings = append(bindings, km.ToggleComplete)
	bindings = append(bindings, km.NextBoard)
	bindings = append(bindings, km.PreviousBoard)
	return bindings
}
