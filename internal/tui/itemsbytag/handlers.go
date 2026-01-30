package itemsbytag

import (
	"fmt"

	"github.com/rhajizada/donezo/internal/tui/styles"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/rhajizada/donezo/internal/tui/navigation"
)

// HandleWindowSize processes window size messages.
func (m *MenuModel) HandleWindowSize(msg tea.WindowSizeMsg) tea.Cmd {
	h, v := styles.App.GetFrameSize()
	m.List.SetSize(msg.Width-h, msg.Height-v)
	return nil
}

// HandleError  processes errors and displays error messages.
func (m *MenuModel) HandleError(msg ErrorMsg) tea.Cmd {
	formattedMsg := fmt.Sprintf("error: %v", msg.Error)
	return m.List.NewStatusMessage(
		styles.ErrorMessage.Render(formattedMsg),
	)
}

// HandleDeleteItem handles DeleteItemMsg.
func (m *MenuModel) HandleDeleteItem(msg DeleteItemMsg) tea.Cmd {
	if msg.Error != nil {
		return m.List.NewStatusMessage(
			styles.ErrorMessage.Render(
				fmt.Sprintf("failed deleting item: %v", msg.Error),
			),
		)
	}

	return m.List.NewStatusMessage(
		styles.StatusMessage.Render(
			fmt.Sprintf("deleted item \"%s\"", msg.Item.Title),
		),
	)
}

func (m *MenuModel) HandleRenameItem(msg RenameItemMsg) tea.Cmd {
	if msg.Error != nil {
		return m.List.NewStatusMessage(
			styles.ErrorMessage.Render(
				fmt.Sprintf("failed renaming item: %v", msg.Error),
			),
		)
	}

	m.List.SetItem(m.List.Index(), NewItem(msg.Item))
	return m.List.NewStatusMessage(
		styles.StatusMessage.Render(
			fmt.Sprintf("updated item \"%s\"", msg.Item.Title),
		),
	)
}

func (m *MenuModel) HandleUpdateTags(msg UpdateTagsMsg) tea.Cmd {
	if msg.Error != nil {
		return m.List.NewStatusMessage(
			styles.ErrorMessage.Render(
				fmt.Sprintf("failed updating tags: %v", msg.Error),
			),
		)
	}

	m.List.SetItem(m.List.Index(), NewItem(msg.Item))
	return m.List.NewStatusMessage(
		styles.StatusMessage.Render(
			fmt.Sprintf("updated item \"%s\" tags", msg.Item.Title),
		),
	)
}

func (m *MenuModel) HandleToggleItem(msg ToggleItemMsg) tea.Cmd {
	if msg.Error != nil {
		return m.List.NewStatusMessage(
			styles.ErrorMessage.Render(
				fmt.Sprintf("failed toggling item: %v", msg.Error),
			),
		)
	}

	m.List.SetItem(m.List.Index(), NewItem(msg.Item))
	selected, ok := m.selectedItem()
	if !ok {
		return m.List.NewStatusMessage(styles.ErrorMessage.Render("no item selected"))
	}
	prefix := ""
	if !selected.Itm.Completed {
		prefix = "in"
	}
	mark := fmt.Sprintf("%scomplete", prefix)

	return m.List.NewStatusMessage(
		styles.StatusMessage.Render(
			fmt.Sprintf("marked item \"%s\" as %s", msg.Item.Title, mark),
		),
	)
}

// HandleInputState handles CreateItemState and RenameItemState states.
func (m *MenuModel) HandleInputState(msg tea.Msg) (textinput.Model, []tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.Input, cmd = m.Input.Update(msg)
	cmds = append(cmds, cmd)

	// Only handle key messages in input states
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		//nolint:exhaustive // Only a subset of key types is handled.
		switch keyMsg.Type {
		case tea.KeyEnter:
			switch m.Context.State {
			case RenameItemNameState:
				m.Context.Title = m.Input.Value()
				m.Context.State = RenameItemDescState
				selected, selectedOK := m.selectedItem()

				if selectedOK {
					m.Input.Placeholder = selected.Itm.Description
					m.Input.SetValue(selected.Itm.Description)
					m.Input.Focus()
				}
			case RenameItemDescState:
				m.Context.Desc = m.Input.Value()
				m.Context.State = DefaultState
				m.Input.Blur()
				cmds = append(cmds, m.RenameItem())
			case UpdateTagsState:
				m.Context.Title = m.Input.Value()
				m.Context.State = DefaultState
				m.Input.Blur()
				cmds = append(cmds, m.UpdateTags())
			case DefaultState:
				// no-op
			default:
				// ignore other states
			}
		case tea.KeyEsc:
			// Cancel the current operation
			m.Context.State = DefaultState
			m.Input.Blur()
		default:
			// ignore other key types
		}
	}

	return m.Input, cmds
}

// HandleKeyInput processes key inputs not handles by list.Model.
func (m *MenuModel) HandleKeyInput(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	if m.List.SettingFilter() {
		return nil
	}
	switch {
	case key.Matches(msg, m.Keys.DeleteItem):
		cmd = m.DeleteItem()
	case key.Matches(msg, m.Keys.RenameItem):
		cmd = m.InitRenameItem()
	case key.Matches(msg, m.Keys.UpdateTags):
		cmd = m.InitUpdateTags()
	case key.Matches(msg, m.Keys.ToggleComplete):
		cmd = m.ToggleComplete()
	case key.Matches(msg, m.Keys.RefreshList):
		cmd = m.ListItems()
	case key.Matches(msg, m.Keys.Back):
		cmd = func() tea.Msg { return navigation.BackMsg{} }
	case key.Matches(msg, m.Keys.NextBoard):
		cmd = func() tea.Msg { return navigation.TagDeltaMsg{Delta: 1} }
	case key.Matches(msg, m.Keys.PreviousBoard):
		cmd = func() tea.Msg { return navigation.TagDeltaMsg{Delta: -1} }
	}
	return cmd
}
