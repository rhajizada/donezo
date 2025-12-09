package itemsbytag

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/tui/helpers"
	"github.com/rhajizada/donezo/internal/tui/styles"
	"github.com/rhajizada/donezo/internal/tui/tags"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *MenuModel) selectedItem() (Item, bool) {
	item, ok := m.List.SelectedItem().(Item)
	return item, ok
}

func (m *MenuModel) selectedTag() (tags.Item, bool) {
	if m.Parent == nil {
		return tags.Item{}, false
	}

	item, ok := m.Parent.List.SelectedItem().(tags.Item)
	return item, ok
}

// ListItems fetches items in the selected board.
func (m *MenuModel) ListItems() tea.Cmd {
	return func() tea.Msg {
		parentItem, ok := m.selectedTag()
		if !ok {
			return ErrorMsg{errors.New("no tag selected")}
		}

		items, err := m.Service.ListItemsByTag(m.ctx, parentItem.Tag)
		if err != nil {
			return ErrorMsg{err}
		}
		return ListItemsMsg{
			items,
		}
	}
}

// RenameItem renames selected item.
func (m *MenuModel) RenameItem() tea.Cmd {
	return func() tea.Msg {
		selected, ok := m.selectedItem()
		if !ok {
			return RenameItemMsg{Error: errors.New("no item selected")}
		}

		selected.Itm.Title = m.Context.Title
		selected.Itm.Description = m.Context.Desc
		item, err := m.Service.UpdateItem(m.ctx, &selected.Itm)
		return RenameItemMsg{
			item,
			err,
		}
	}
}

// UpdateTags updates item tags.
func (m *MenuModel) UpdateTags() tea.Cmd {
	return func() tea.Msg {
		var item *service.Item
		selected, ok := m.selectedItem()
		if !ok {
			return UpdateTagsMsg{Item: nil, Error: errors.New("no item selected")}
		}
		tags, err := helpers.ExtractTags(m.Context.Title)
		if err != nil {
			return UpdateTagsMsg{
				item,
				err,
			}
		}

		selected.Itm.Tags = tags
		item, err = m.Service.UpdateItem(m.ctx, &selected.Itm)
		return UpdateTagsMsg{
			item,
			err,
		}
	}
}

// InitRenameItem starts the renaming process for the selected item.
func (m *MenuModel) InitRenameItem() tea.Cmd {
	if len(m.List.Items()) == 0 {
		return m.List.NewStatusMessage(
			styles.StatusMessage.Render("no item selected"))
	}

	m.Context.State = RenameItemNameState
	selected, ok := m.selectedItem()
	if ok {
		m.Input.SetValue(selected.Itm.Title)
	}
	m.Input.Focus()
	return nil
}

// InitUpdateTags initializes tag updates.
func (m *MenuModel) InitUpdateTags() tea.Cmd {
	m.Context.State = UpdateTagsState
	m.Input.Placeholder = "Enter comma-separated list of tags"
	selected, ok := m.selectedItem()
	if ok {
		dSep := fmt.Sprintf(" %s", helpers.TagsSeparator)
		m.Input.SetValue(strings.Join(selected.Itm.Tags, dSep))
	}
	m.Input.Focus()
	return nil
}

// DeleteItem deletes current selected item.
func (m *MenuModel) DeleteItem() tea.Cmd {
	return func() tea.Msg {
		selected, ok := m.selectedItem()
		if !ok {
			return DeleteItemMsg{Error: errors.New("no item selected")}
		}

		err := m.Service.DeleteItem(m.ctx, &selected.Itm)
		return DeleteItemMsg{Error: err, Item: &selected.Itm}
	}
}

func (m MenuModel) ToggleComplete() tea.Cmd {
	if len(m.List.Items()) == 0 {
		return m.List.NewStatusMessage(
			styles.ErrorMessage.Render("no item selected"))
	}

	selected, ok := m.selectedItem()
	if !ok {
		return m.List.NewStatusMessage(styles.ErrorMessage.Render("no item selected"))
	}
	selected.Itm.Completed = !selected.Itm.Completed
	m.List.SetItem(m.List.Index(), selected)

	return func() tea.Msg {
		i, err := m.Service.UpdateItem(m.ctx, &selected.Itm)
		return ToggleItemMsg{i, err}
	}
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if m.Context.State != DefaultState {
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

	case ListItemsMsg:
		m.List.SetItems(NewList(msg.Items))

	case DeleteItemMsg:
		cmd := m.HandleDeleteItem(msg)
		cmds = append(cmds, cmd)
		cmd = m.ListItems()
		cmds = append(cmds, cmd)

	case RenameItemMsg:
		cmd := m.HandleRenameItem(msg)
		cmds = append(cmds, cmd)

	case UpdateTagsMsg:
		cmd := m.HandleUpdateTags(msg)
		cmds = append(cmds, cmd)

	case ToggleItemMsg:
		cmd := m.HandleToggleItem(msg)
		cmds = append(cmds, cmd)
	}

	listModel, listCmd := m.List.Update(msg)
	m.List = listModel
	cmds = append(cmds, listCmd)

	return m, tea.Batch(cmds...)
}
