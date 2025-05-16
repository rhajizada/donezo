package tags

import (
	"fmt"

	"github.com/rhajizada/donezo/internal/tui/styles"

	"github.com/charmbracelet/bubbles/key"
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

// HandleDeleteTag handles DeleteTagMsg
func (m *MenuModel) HandleDeleteTag(msg DeleteTagMsg) tea.Cmd {
	if msg.Error != nil {
		return m.List.NewStatusMessage(
			styles.ErrorMessage.Render(
				fmt.Sprintf("failed deleting tag: %v", msg.Error),
			),
		)
	} else {
		return m.List.NewStatusMessage(
			styles.StatusMessage.Render(
				fmt.Sprintf("deleted tag \"%s\"", msg.Tag),
			),
		)
	}
}

// HandleKeyInput processes key inputs not handles by list.Model
func (m *MenuModel) HandleKeyInput(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch {
	case key.Matches(msg, m.Keys.DeleteTag):
		cmd = m.DeleteTag()
	case key.Matches(msg, m.Keys.RefreshList):
		cmd = m.ListTags()
	case key.Matches(msg, m.Keys.Copy):
		cmd = m.Copy()
	}
	return cmd
}
