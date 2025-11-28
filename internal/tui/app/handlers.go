package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/rhajizada/donezo/internal/tui/boards"
	"github.com/rhajizada/donezo/internal/tui/itemsbyboard"
	"github.com/rhajizada/donezo/internal/tui/itemsbytag"
	"github.com/rhajizada/donezo/internal/tui/navigation"
)

func (m AppModel) switchMain(view navigation.View) (tea.Model, tea.Cmd) {
	if view == m.active {
		return m, nil
	}

	switch view {
	case navigation.ViewBoards:
		m.active = navigation.ViewBoards
		return m, m.initWithSize(m.boards.Init())
	case navigation.ViewTags:
		m.active = navigation.ViewTags
		return m, m.initWithSize(m.tags.Init())
	default:
		return m, nil
	}
}

func (m AppModel) openBoardItems() (tea.Model, tea.Cmd) {
	if m.boards == nil || m.boards.List.SettingFilter() || m.boards.State != boards.DefaultState {
		return m, nil
	}
	if len(m.boards.List.Items()) == 0 {
		return m, nil
	}
	itemMenu := itemsbyboard.New(m.ctx, m.service, m.boards)
	m.itemsByBoard = &itemMenu
	m.active = navigation.ViewItemsByBoard
	return m, m.initWithSize(itemMenu.Init())
}

func (m AppModel) openTagItems() (tea.Model, tea.Cmd) {
	if m.tags == nil || m.tags.List.SettingFilter() {
		return m, nil
	}
	if len(m.tags.List.Items()) == 0 {
		return m, nil
	}
	itemMenu := itemsbytag.New(m.ctx, m.service, m.tags)
	m.itemsByTag = &itemMenu
	m.active = navigation.ViewItemsByTag
	return m, m.initWithSize(itemMenu.Init())
}

func (m AppModel) navigateBack() (tea.Model, tea.Cmd) {
	switch m.active {
	case navigation.ViewItemsByBoard:
		m.active = navigation.ViewBoards
		return m, m.forwardCachedSize()
	case navigation.ViewItemsByTag:
		m.active = navigation.ViewTags
		return m, m.forwardCachedSize()
	default:
		return m, nil
	}
}

func (m AppModel) moveBoardSelection(delta int) (tea.Model, tea.Cmd) {
	if m.boards == nil || m.boards.State != boards.DefaultState || m.boards.List.SettingFilter() {
		return m, nil
	}
	items := m.boards.List.Items()
	if len(items) == 0 {
		return m, nil
	}

	current := m.boards.List.Index()
	nextIndex := (current + delta) % len(items)
	if nextIndex < 0 {
		nextIndex += len(items)
	}

	m.boards.List.Select(nextIndex)
	if m.active == navigation.ViewItemsByBoard {
		return m.openBoardItems()
	}
	return m, nil
}

func (m AppModel) moveTagSelection(delta int) (tea.Model, tea.Cmd) {
	if m.tags == nil || m.tags.List.SettingFilter() {
		return m, nil
	}
	items := m.tags.List.Items()
	if len(items) == 0 {
		return m, nil
	}

	current := m.tags.List.Index()
	nextIndex := (current + delta) % len(items)
	if nextIndex < 0 {
		nextIndex += len(items)
	}

	m.tags.List.Select(nextIndex)
	if m.active == navigation.ViewItemsByTag {
		return m.openTagItems()
	}
	return m, nil
}

func (m AppModel) initWithSize(initCmd tea.Cmd) tea.Cmd {
	if m.lastSize == nil {
		return initCmd
	}
	size := *m.lastSize
	return tea.Batch(initCmd, func() tea.Msg { return size })
}

func (m AppModel) forwardCachedSize() tea.Cmd {
	if m.lastSize == nil {
		return nil
	}
	size := *m.lastSize
	return func() tea.Msg { return size }
}
