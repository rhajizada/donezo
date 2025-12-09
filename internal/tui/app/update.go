package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/rhajizada/donezo/internal/tui/boards"
	"github.com/rhajizada/donezo/internal/tui/itemsbyboard"
	"github.com/rhajizada/donezo/internal/tui/itemsbytag"
	"github.com/rhajizada/donezo/internal/tui/navigation"
	"github.com/rhajizada/donezo/internal/tui/tags"
)

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.lastSize = &msg
		return m, m.forwardWindowSize(msg)
	case navigation.SwitchMainViewMsg:
		return m.switchMain(msg.View)
	case navigation.OpenBoardItemsMsg:
		return m.openBoardItems()
	case navigation.OpenTagItemsMsg:
		return m.openTagItems()
	case navigation.BackMsg:
		return m.navigateBack()
	case navigation.BoardDeltaMsg:
		return m.moveBoardSelection(msg.Delta)
	case navigation.TagDeltaMsg:
		return m.moveTagSelection(msg.Delta)
	}

	active := m.activeModel()
	if active == nil {
		return m, nil
	}

	updated, cmd := active.Update(msg)
	m.setActiveModel(updated)
	return m, cmd
}

func (m *AppModel) forwardWindowSize(msg tea.WindowSizeMsg) tea.Cmd {
	active := m.activeModel()
	if active == nil {
		return nil
	}
	updated, cmd := active.Update(msg)
	m.setActiveModel(updated)
	return cmd
}

func (m *AppModel) setActiveModel(model tea.Model) {
	switch m.active {
	case navigation.ViewBoards:
		switch v := model.(type) {
		case boards.MenuModel:
			m.boards = &v
		case *boards.MenuModel:
			m.boards = v
		}
	case navigation.ViewTags:
		switch v := model.(type) {
		case tags.MenuModel:
			m.tags = &v
		case *tags.MenuModel:
			m.tags = v
		}
	case navigation.ViewItemsByBoard:
		switch v := model.(type) {
		case itemsbyboard.MenuModel:
			m.itemsByBoard = &v
		case *itemsbyboard.MenuModel:
			m.itemsByBoard = v
		}
	case navigation.ViewItemsByTag:
		switch v := model.(type) {
		case itemsbytag.MenuModel:
			m.itemsByTag = &v
		case *itemsbytag.MenuModel:
			m.itemsByTag = v
		}
	}
}

func (m *AppModel) activeModel() tea.Model {
	switch m.active {
	case navigation.ViewBoards:
		return m.boards
	case navigation.ViewTags:
		return m.tags
	case navigation.ViewItemsByBoard:
		if m.itemsByBoard != nil {
			return m.itemsByBoard
		}
	case navigation.ViewItemsByTag:
		if m.itemsByTag != nil {
			return m.itemsByTag
		}
	}
	return nil
}
