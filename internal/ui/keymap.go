package ui

import "github.com/charmbracelet/bubbles/key"

// Keymap defines custom key bindings and implements help.KeyMap interface.
type Keymap struct {
	AddItem          key.Binding
	DeleteItem       key.Binding
	RenameItem       key.Binding
	RefreshList      key.Binding
	NextBoard        key.Binding
	PrevBoard        key.Binding
	ToggleTitleBar   key.Binding
	ToggleStatusBar  key.Binding
	TogglePagination key.Binding
	ToggleComplete   key.Binding
	ToggleHelpMenu   key.Binding
	Quit             key.Binding
}

// newKeymap initializes a new listKeyMap with custom bindings.
func newKeymap() *Keymap {
	return &Keymap{
		AddItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
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
		NextBoard: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next board"),
		),
		PrevBoard: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev board"),
		),
		ToggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		ToggleStatusBar: key.NewBinding(
			key.WithKeys("B"),
			key.WithHelp("B", "toggle status"),
		),
		TogglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		ToggleComplete: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle complete"),
		),
		ToggleHelpMenu: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q", "quit"),
		),
	}
}

// ShortHelp returns keybindings to be shown in the mini help view.
// It's part of the help.KeyMap interface.
func (k Keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.AddItem,
		k.NextBoard,
	}
}

// FullHelp returns keybindings for the expanded help view.
// It's part of the help.KeyMap interface.
func (k Keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.AddItem, k.DeleteItem, k.RenameItem, k.RefreshList},
		{k.NextBoard, k.PrevBoard, k.ToggleTitleBar, k.ToggleStatusBar},
		{k.ToggleComplete},
	}
}
