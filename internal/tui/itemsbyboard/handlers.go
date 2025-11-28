package itemsbyboard

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

// HandleError  processes errors and displays error messages
func (m *MenuModel) HandleError(msg ErrorMsg) tea.Cmd {
	formattedMsg := fmt.Sprintf("error: %v", msg.Error)
	return m.List.NewStatusMessage(
		styles.ErrorMessage.Render(formattedMsg),
	)
}

// HandleCreateItem handles CreateItemMsg
func (m *MenuModel) HandleCreateItem(msg CreateItemMsg) tea.Cmd {
	if msg.Error != nil {
		return m.List.NewStatusMessage(
			styles.ErrorMessage.Render(
				fmt.Sprintf("error creating item: %v", msg.Error),
			),
		)
	}
	m.List.InsertItem(len(m.List.Items()), NewItem(msg.Item))
	return m.List.NewStatusMessage(
		styles.StatusMessage.Render(
			fmt.Sprintf("created item \"%s\"", msg.Item.Title),
		),
	)
}

// HandleDeleteItem handles DeleteItemMsg
func (m *MenuModel) HandleDeleteItem(msg DeleteItemMsg) tea.Cmd {
	if msg.Error != nil {
		return m.List.NewStatusMessage(
			styles.ErrorMessage.Render(
				fmt.Sprintf("failed deleting item: %v", msg.Error),
			),
		)
	} else {
		return m.List.NewStatusMessage(
			styles.StatusMessage.Render(
				fmt.Sprintf("deleted item \"%s\"", msg.Item.Title),
			),
		)
	}
}

func (m *MenuModel) HandleRenameItem(msg RenameItemMsg) tea.Cmd {
	if msg.Error != nil {
		return m.List.NewStatusMessage(
			styles.ErrorMessage.Render(
				fmt.Sprintf("failed renaming item: %v", msg.Error),
			),
		)
	} else {
		m.List.SetItem(m.List.Index(), NewItem(msg.Item))
		return m.List.NewStatusMessage(
			styles.StatusMessage.Render(
				fmt.Sprintf("updated item \"%s\"", msg.Item.Title),
			),
		)
	}
}

func (m *MenuModel) HandleUpdateTags(msg UpdateTagsMsg) tea.Cmd {
	if msg.Error != nil {
		return m.List.NewStatusMessage(
			styles.ErrorMessage.Render(
				fmt.Sprintf("failed updating tags: %v", msg.Error),
			),
		)
	} else {
		m.List.SetItem(m.List.Index(), NewItem(msg.Item))
		return m.List.NewStatusMessage(
			styles.StatusMessage.Render(
				fmt.Sprintf("updated item \"%s\" tags", msg.Item.Title),
			),
		)
	}
}

func (m *MenuModel) HandleToggleItem(msg ToggleItemMsg) tea.Cmd {
	if msg.Error != nil {
		return m.List.NewStatusMessage(
			styles.ErrorMessage.Render(
				fmt.Sprintf("failed toggling item: %v", msg.Error),
			),
		)
	} else {
		m.List.SetItem(m.List.Index(), NewItem(msg.Item))
		selected := m.List.SelectedItem().(Item)
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
}

// HandleInputState handles CreateItemState and RenameItemState states
func (m *MenuModel) HandleInputState(msg tea.Msg) (textinput.Model, []tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.Input, cmd = m.Input.Update(msg)
	cmds = append(cmds, cmd)

	// Only handle key messages in input states
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyEnter:
			switch m.Context.State {
			case CreateItemNameState:
				m.Context.Title = m.Input.Value()
				m.Context.State = CreateItemDescState
				m.Input.Placeholder = "Enter item description"
				m.Input.SetValue("")
				m.Input.Focus()
			case CreateItemDescState:
				m.Context.Desc = m.Input.Value()
				m.Context.State = DefaultState
				m.Input.Blur()
				cmds = append(cmds, m.CreateItem())
			case RenameItemNameState:
				m.Context.Title = m.Input.Value()
				m.Context.State = RenameItemDescState
				selected, ok := m.List.SelectedItem().(Item)

				if ok {
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
			}
		case tea.KeyEsc:
			// Cancel the current operation
			m.Context.State = DefaultState
			m.Input.Blur()
		}
	}

	return m.Input, cmds
}

// HandleKeyInput processes key inputs not handles by list.Model
func (m *MenuModel) HandleKeyInput(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	if m.List.SettingFilter() {
		return nil
	}
	switch {
	case key.Matches(msg, m.Keys.CreateItem):
		cmd = m.InitCreateItem()
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
	case key.Matches(msg, m.Keys.Copy):
		cmd = m.Copy()
	case key.Matches(msg, m.Keys.Paste):
		cmd = m.Paste()
	case key.Matches(msg, m.Keys.Back):
		cmd = func() tea.Msg { return navigation.BackMsg{} }
	case key.Matches(msg, m.Keys.NextBoard):
		cmd = func() tea.Msg { return navigation.BoardDeltaMsg{Delta: 1} }
	case key.Matches(msg, m.Keys.PreviousBoard):
		cmd = func() tea.Msg { return navigation.BoardDeltaMsg{Delta: -1} }
	}
	return cmd
}
