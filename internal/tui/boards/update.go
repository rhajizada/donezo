package boards

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// ListBoards fetches the list of boards from the client.
func (m *MenuModel) ListBoards() tea.Cmd {
	return func() tea.Msg {
		boards, err := m.Client.ListBoards()
		if err != nil {
			return ErrorMsg{err}
		}
		if len(*boards) == 0 {
			return ErrorMsg{fmt.Errorf("no boards available")}
		}
		return ListBoardsMsg{
			boards,
		}
	}
}

// CreateBoard creates a new board
func (m *MenuModel) CreateBoard() tea.Cmd {
	return func() tea.Msg {
		board, err := m.Client.CreateBoard(m.Input.Value())
		return CreateBoardMsg{
			board,
			err,
		}
	}
}

// RenameBoard renames selected board
func (m *MenuModel) RenameBoard() tea.Cmd {
	return func() tea.Msg {
		selected := m.List.SelectedItem().(Item)
		selected.Board.Name = m.Input.Value()
		board, err := m.Client.UpdateBoard(&selected.Board)
		return RenameBoardMsg{
			board,
			err,
		}
	}
}

// InitCreateBoard sets list state to CreateBoardState to render text input
func (m *MenuModel) InitCreateBoard() tea.Cmd {
	m.State = CreateBoardState
	m.Input.Placeholder = "Enter board name"
	m.Input.SetValue("")
	m.Input.Focus()
	return nil
}

// DeleteBoard deletes current selected board
func (m *MenuModel) DeleteBoard() tea.Cmd {
	return func() tea.Msg {
		selected := m.List.SelectedItem().(Item)
		err := m.Client.DeleteBoard(&selected.Board)
		return DeleteBoardMsg{Error: err, Board: &selected.Board}
	}
}

// InitRenameBoard sets list state to InitRenameBoard to render text input
func (m *MenuModel) InitRenameBoard() tea.Cmd {
	m.State = RenameBoardState
	selected := m.List.SelectedItem().(Item)
	m.Input.SetValue(selected.Board.Name)
	m.Input.Focus()
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if m.State != DefaultState {
		m.Input, cmds = m.HandleInputState(msg)
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cmd := m.HandleWindowSize(msg)
		cmds = append(cmds, cmd)

	case tea.KeyMsg:
		cmd := m.HandleKeyInput(msg)
		cmds = append(cmds, cmd)

	case ErrorMsg:
		cmd := m.HandleError(msg)
		cmds = append(cmds, cmd)

	case ListBoardsMsg:
		m.List.SetItems(NewList(msg.Boards))

	case CreateBoardMsg:
		cmd := m.HandleCreateBoard(msg)
		cmds = append(cmds, cmd)

	case DeleteBoardMsg:
		cmd := m.HandleDeleteBoard(msg)
		cmds = append(cmds, cmd)
		cmd = m.ListBoards()
		cmds = append(cmds, cmd)

	case RenameBoardMsg:
		cmd := m.HandleRenameBoard(msg)
		cmds = append(cmds, cmd)

	}

	listModel, listCmd := m.List.Update(msg)
	m.List = listModel
	cmds = append(cmds, listCmd)

	return m, tea.Batch(cmds...)
}
