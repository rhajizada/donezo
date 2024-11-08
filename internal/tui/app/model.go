package app

import (
	"github.com/rhajizada/donezo/internal/tui/boards"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rhajizada/donezo/pkg/client"
)

type AppModel struct {
	Client         *client.Client
	ViewStack      []tea.Model
	LastWindowSize *tea.WindowSizeMsg // Store the last known WindowSizeMsg
}

func NewModel(client *client.Client) AppModel {
	boardMenu := boards.NewModel(client)
	return AppModel{
		Client:         client,
		ViewStack:      []tea.Model{boardMenu},
		LastWindowSize: nil, // Initialize without a WindowSizeMsg
	}
}

func (m AppModel) Init() tea.Cmd {
	// Initialize the top model in the stack
	return m.ViewStack[len(m.ViewStack)-1].Init()
}

func (m *AppModel) GetCurrentBoard() *client.Board {
	// Retrieve the currently selected board from the board menu
	if boardMenu, ok := m.ViewStack[0].(boards.MenuModel); ok {
		if selected, ok := boardMenu.List.SelectedItem().(boards.Item); ok {
			return &selected.Board
		}
	}
	return nil
}
