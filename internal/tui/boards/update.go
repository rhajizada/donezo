package boards

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/tui/styles"
	"golang.design/x/clipboard"
)

var writeClipboardText = func(data []byte) {
	clipboard.Write(clipboard.FmtText, data)
}

func (m *MenuModel) selectedItem() (Item, bool) {
	item, ok := m.List.SelectedItem().(Item)
	return item, ok
}

// ListBoards fetches the list of boards from the client.
func (m *MenuModel) ListBoards() tea.Cmd {
	return func() tea.Msg {
		boards, err := m.Client.ListBoards(m.ctx)
		if err != nil {
			return ErrorMsg{err}
		}
		return ListBoardsMsg{
			boards,
		}
	}
}

// Copy copies board name to system clipboard.
func (m *MenuModel) Copy() tea.Cmd {
	selected, ok := m.selectedItem()
	if !ok {
		return nil
	}
	currentBoard := selected.Board
	items, err := m.Client.ListItemsByBoard(m.ctx, &currentBoard)
	if err != nil {
		return func() tea.Msg {
			return ErrorMsg{err}
		}
	}
	md := service.ItemsToMarkdown(currentBoard.Name, *items)
	writeClipboardText([]byte(md))
	return m.List.NewStatusMessage(
		styles.StatusMessage.Render(
			fmt.Sprintf("copied \"%s\" to system clipboard", currentBoard.Name),
		),
	)
}

// CreateBoard creates a new board.
func (m *MenuModel) CreateBoard() tea.Cmd {
	return func() tea.Msg {
		board, err := m.Client.CreateBoard(m.ctx, m.Input.Value())
		return CreateBoardMsg{
			board,
			err,
		}
	}
}

// RenameBoard renames selected board.
func (m *MenuModel) RenameBoard() tea.Cmd {
	return func() tea.Msg {
		selected, ok := m.selectedItem()
		if !ok {
			return ErrorMsg{errors.New("no board selected")}
		}
		selected.Board.Name = m.Input.Value()
		board, err := m.Client.UpdateBoard(m.ctx, &selected.Board)
		return RenameBoardMsg{
			board,
			err,
		}
	}
}

// InitCreateBoard sets list state to CreateBoardState to render text input.
func (m *MenuModel) InitCreateBoard() tea.Cmd {
	m.State = CreateBoardState
	m.Input.Placeholder = "Enter board name"
	m.Input.SetValue("")
	m.Input.Focus()
	return nil
}

// DeleteBoard deletes current selected board.
func (m *MenuModel) DeleteBoard() tea.Cmd {
	return func() tea.Msg {
		selected, ok := m.selectedItem()
		if !ok {
			return DeleteBoardMsg{Error: errors.New("no board selected")}
		}
		err := m.Client.DeleteBoard(m.ctx, &selected.Board)
		return DeleteBoardMsg{Error: err, Board: &selected.Board}
	}
}

// InitRenameBoard sets list state to InitRenameBoard to render text input.
func (m *MenuModel) InitRenameBoard() tea.Cmd {
	m.State = RenameBoardState
	selected, ok := m.selectedItem()
	if ok {
		m.Input.SetValue(selected.Board.Name)
	}
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

	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEsc {
		return m, tea.Batch(cmds...)
	}

	listModel, listCmd := m.List.Update(msg)
	m.List = listModel
	cmds = append(cmds, listCmd)

	return m, tea.Batch(cmds...)
}
