package app

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/tui/boards"
	"github.com/rhajizada/donezo/internal/tui/itemsbyboard"
	"github.com/rhajizada/donezo/internal/tui/itemsbytag"
	"github.com/rhajizada/donezo/internal/tui/navigation"
	"github.com/rhajizada/donezo/internal/tui/tags"
)

//revive:disable-next-line:exported // app.AppModel is the public entry point for the TUI.
//nolint:recvcheck // Mixed receivers align with tea.Model usage patterns.
type AppModel struct {
	ctx     context.Context
	service *service.Service

	boards       *boards.MenuModel
	tags         *tags.MenuModel
	itemsByBoard *itemsbyboard.MenuModel
	itemsByTag   *itemsbytag.MenuModel

	active   navigation.View
	lastSize *tea.WindowSizeMsg
}

func New(ctx context.Context, service *service.Service) AppModel {
	boardMenu := boards.New(ctx, service)
	tagMenu := tags.NewModel(ctx, service)
	return AppModel{
		ctx:     ctx,
		service: service,
		boards:  &boardMenu,
		tags:    &tagMenu,
		active:  navigation.ViewBoards,
	}
}

func (m AppModel) Init() tea.Cmd {
	return m.activeModel().Init()
}
