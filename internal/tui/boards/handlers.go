package boards

import (
	"fmt"

	"github.com/rhajizada/donezo/internal/tui/styles"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// HandleWindowSize processes window size messages.
func (m *MenuModel) HandleWindowSize(msg tea.WindowSizeMsg) tea.Cmd {
	h, v := styles.App.GetFrameSize()
	m.List.SetSize(msg.Width-h, msg.Height-v)
	return nil
}

// HandleError  processes errors and displays error messages
func (m *MenuModel) HandleError(msg ErrorMsg) tea.Cmd {
	formattedMsg := fmt.Sprintf("error: %v", msg.Error)
	return m.List.NewStatusMessage(
		styles.ErrorMessage.Render(formattedMsg),
	)
}

// HandleCreateBoard handles CreateBoardMsg
func (m *MenuModel) HandleCreateBoard(msg CreateBoardMsg) tea.Cmd {
	if msg.Error != nil {
		return m.List.NewStatusMessage(
			styles.ErrorMessage.Render(
				fmt.Sprintf("error creating board: %v", msg.Error),
			),
		)
	}
	m.List.InsertItem(len(m.List.Items()), Item{Board: *msg.Board})
	return m.List.NewStatusMessage(
		styles.StatusMessage.Render(
			fmt.Sprintf("created board \"%s\"", msg.Board.Name),
		),
	)
}

// HandleDeleteBoard handles DeleteBoardMsg
func (m *MenuModel) HandleDeleteBoard(msg DeleteBoardMsg) tea.Cmd {
	if msg.Error != nil {
		return m.List.NewStatusMessage(
			styles.ErrorMessage.Render(
				fmt.Sprintf("failed deleting board: %v", msg.Error),
			),
		)
	} else {
		return m.List.NewStatusMessage(
			styles.StatusMessage.Render(
				fmt.Sprintf("deleted board \"%s\"", msg.Board.Name),
			),
		)
	}
}

func (m *MenuModel) HandleRenameBoard(msg RenameBoardMsg) tea.Cmd {
	if msg.Error != nil {
		return m.List.NewStatusMessage(
			styles.ErrorMessage.Render(
				fmt.Sprintf("failed renaming board: %v", msg.Error),
			),
		)
	} else {
		m.List.SetItem(m.List.Index(), NewItem(msg.Board))
		return m.List.NewStatusMessage(
			styles.StatusMessage.Render(
				fmt.Sprintf("renamed board to \"%s\"", msg.Board.Name),
			),
		)
	}
}

// HandleInputState handles CreateBoardState and RenameBoardState states
func (m *MenuModel) HandleInputState(msg tea.Msg) (textinput.Model, []tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.Input, cmd = m.Input.Update(msg)
	cmds = append(cmds, cmd)

	// Only handle key messages in input states
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyEnter:
			switch m.State {
			case CreateBoardState:
				cmds = append(cmds, m.CreateBoard())
				m.State = DefaultState
				m.Input.Blur()
			case RenameBoardState:
				cmds = append(cmds, m.RenameBoard())
				m.State = DefaultState
				m.Input.Blur()
			}
		case tea.KeyEsc:
			// Cancel the current operation
			m.State = DefaultState
			m.Input.Blur()
		}
	}

	return m.Input, cmds
}

// HandleKeyInput processes key inputs not handles by list.Model
func (m *MenuModel) HandleKeyInput(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	if !m.List.SettingFilter() {
		switch {
		case key.Matches(msg, m.Keys.CreateBoard):
			cmd = m.InitCreateBoard()
		case key.Matches(msg, m.Keys.DeleteBoard):
			cmd = m.DeleteBoard()
		case key.Matches(msg, m.Keys.RenameBoard):
			cmd = m.InitRenameBoard()
		case key.Matches(msg, m.Keys.RefreshList):
			cmd = m.ListBoards()
		case key.Matches(msg, m.Keys.Copy):
			cmd = m.Copy()
		}
	}
	return cmd
}
