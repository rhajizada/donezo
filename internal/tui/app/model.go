package app

import (
	"context"

	"github.com/rhajizada/donezo/internal/tui/boards"
	"github.com/rhajizada/donezo/internal/tui/tags"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rhajizada/donezo/internal/service"
)

type AppModel struct {
	ctx            context.Context
	Service        *service.Service
	ViewStack      []tea.Model
	LastWindowSize *tea.WindowSizeMsg // Store the last known WindowSizeMsg
	MenuType       MenuType           // NEW: tracks whether main view is boards or tags
}

func New(ctx context.Context, service *service.Service) AppModel {
	boardMenu := boards.New(ctx, service)
	return AppModel{
		ctx:            ctx,
		Service:        service,
		ViewStack:      []tea.Model{boardMenu},
		LastWindowSize: nil,
		MenuType:       MenuBoards, // initially showing boards
	}
}

func (m AppModel) Init() tea.Cmd {
	// Initialize the top model in the stack
	return m.ViewStack[len(m.ViewStack)-1].Init()
}

func (m *AppModel) GetCurrentBoard() *boards.MenuModel {
	// Retrieve the currently selected board from the board menu
	if boardMenu, ok := m.ViewStack[0].(boards.MenuModel); ok {
		if _, ok := boardMenu.List.SelectedItem().(boards.Item); ok {
			return &boardMenu
		}
	}
	return nil
}

func (m *AppModel) GetCurrentTag() string {
	if tagMenu, ok := m.ViewStack[0].(tags.MenuModel); ok {
		if selected, ok := tagMenu.List.SelectedItem().(tags.Item); ok {
			return selected.Tag
		}
	}
	return ""
}
