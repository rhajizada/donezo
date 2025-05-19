package itemsbyboard

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/tui/boards"
	"github.com/rhajizada/donezo/internal/tui/styles"
	"golang.design/x/clipboard"

	tea "github.com/charmbracelet/bubbletea"
)

const TagsSeparator = ","

// Copy copies selected item to clipboard and moves it to ItemStack
func (m *MenuModel) Copy() tea.Cmd {
	selectedItem := m.List.SelectedItem().(Item).Itm
	data, err := json.Marshal(selectedItem)
	if err != nil {
		return func() tea.Msg {
			return ErrorMsg{err}
		}
	}

	clipboard.Write(clipboard.FmtText, data)
	return m.List.NewStatusMessage(
		styles.StatusMessage.Render(
			fmt.Sprintf("copied \"%s\" to system clipboard", selectedItem.Title),
		),
	)
}

// Paste pastes item into current board
func (m *MenuModel) Paste() tea.Cmd {
	data := clipboard.Read(clipboard.FmtText)
	var lastItem service.Item
	err := json.Unmarshal(data, &lastItem)
	if err != nil {
		return m.List.NewStatusMessage(
			styles.ErrorMessage.Render(
				fmt.Sprintf("no items in clipboard: %v", err),
			),
		)
	}
	currentBoard := m.Parent.List.SelectedItem().(boards.Item).Board
	item, err := m.Service.CreateItem(m.ctx, &currentBoard, lastItem.Title, lastItem.Description)
	if err != nil {
		return func() tea.Msg {
			return ErrorMsg{err}
		}
	}
	item.Tags = lastItem.Tags
	item.Completed = lastItem.Completed
	item, err = m.Service.UpdateItem(m.ctx, item)
	return func() tea.Msg {
		return CreateItemMsg{item, err}
	}
}

// ListItems fetches items ine the selected board.
func (m *MenuModel) ListItems() tea.Cmd {
	return func() tea.Msg {
		parentItem := m.Parent.List.SelectedItem().(boards.Item)
		items, err := m.Service.ListItemsByBoard(m.ctx, &parentItem.Board)
		if err != nil {
			return ErrorMsg{err}
		}
		return ListItemsMsg{
			items,
		}
	}
}

// CreateItem creates a new item
func (m *MenuModel) CreateItem() tea.Cmd {
	return func() tea.Msg {
		parentItem := m.Parent.List.SelectedItem().(boards.Item)
		item, err := m.Service.CreateItem(m.ctx, &parentItem.Board, m.Context.Title, m.Context.Desc)
		return CreateItemMsg{
			item,
			err,
		}
	}
}

func (m *MenuModel) InitCreateItem() tea.Cmd {
	m.Context.State = CreateItemNameState
	m.Input.Placeholder = "Enter item name"
	m.Input.SetValue("")
	m.Input.Focus()
	return nil
}

// RenameItem renames selected item
func (m *MenuModel) RenameItem() tea.Cmd {
	return func() tea.Msg {
		selected := m.List.SelectedItem().(Item)
		selected.Itm.Title = m.Context.Title
		selected.Itm.Description = m.Context.Desc
		item, err := m.Service.UpdateItem(m.ctx, &selected.Itm)
		return RenameItemMsg{
			item,
			err,
		}
	}
}

// UpdateTags updates item tags
func (m *MenuModel) UpdateTags() tea.Cmd {
	return func() tea.Msg {
		selected := m.List.SelectedItem().(Item)
		tags := strings.Split(m.Context.Title, TagsSeparator)
		emptyTag := false
		for idx, tag := range tags {
			sanitized := strings.TrimSpace(tag)
			if len(sanitized) == 0 {
				emptyTag = true
				break
			} else {
				tags[idx] = sanitized
			}
		}

		var item *service.Item
		var err error
		if !emptyTag {
			selected.Itm.Tags = tags
			item, err = m.Service.UpdateItem(m.ctx, &selected.Itm)
		} else {
			err = errors.New("tag must not be empty")
		}
		return UpdateTagsMsg{
			item,
			err,
		}
	}
}

// initiateRename starts the renaming process for the selected item.
func (m *MenuModel) InitRenameItem() tea.Cmd {
	if len(m.List.Items()) == 0 {
		return m.List.NewStatusMessage(
			styles.StatusMessage.Render("no item selected"))
	}

	m.Context.State = RenameItemNameState
	selected := m.List.SelectedItem().(Item)
	m.Input.SetValue(selected.Itm.Title)
	m.Input.Focus()
	return nil
}

// InitUpdateTags initiaizes tag updates
func (m *MenuModel) InitUpdateTags() tea.Cmd {
	m.Context.State = UpdateTagsState
	m.Input.Placeholder = "Enter comma-separated list of tags"
	selected := m.List.SelectedItem().(Item)
	dSep := fmt.Sprintf(" %s", TagsSeparator)
	m.Input.SetValue(strings.Join(selected.Itm.Tags, dSep))
	m.Input.Focus()
	return nil
}

// DeleteBoard deletes current selected board
func (m *MenuModel) DeleteItem() tea.Cmd {
	m.Copy()
	selectedItem := m.List.SelectedItem().(Item).Itm
	err := m.Service.DeleteItem(m.ctx, &selectedItem)
	return func() tea.Msg {
		return DeleteItemMsg{Error: err, Item: &selectedItem}
	}
}

func (m MenuModel) ToggleComplete() tea.Cmd {
	if len(m.List.Items()) == 0 {
		return m.List.NewStatusMessage(
			styles.ErrorMessage.Render("no item selected"))
	}

	selected := m.List.SelectedItem().(Item)
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

	case CreateItemMsg:
		cmd := m.HandleCreateItem(msg)
		cmds = append(cmds, cmd)

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
