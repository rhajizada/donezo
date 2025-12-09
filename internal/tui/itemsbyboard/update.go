package itemsbyboard

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/tui/boards"
	"github.com/rhajizada/donezo/internal/tui/helpers"
	"github.com/rhajizada/donezo/internal/tui/styles"
	"golang.design/x/clipboard"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	//nolint:gochecknoglobals // injectable for tests
	writeClipboardText = func(data []byte) {
		clipboard.Write(clipboard.FmtText, data)
	}
	//nolint:gochecknoglobals // injectable for tests
	readClipboardText = func() []byte {
		return clipboard.Read(clipboard.FmtText)
	}
)

func (m *MenuModel) selectedItem() (Item, bool) {
	item, ok := m.List.SelectedItem().(Item)
	return item, ok
}

func (m *MenuModel) selectedBoard() (boards.Item, bool) {
	if m.Parent == nil {
		return boards.Item{}, false
	}

	item, ok := m.Parent.List.SelectedItem().(boards.Item)
	return item, ok
}

// Copy copies selected item to clipboard and moves it to ItemStack.
func (m *MenuModel) Copy() tea.Cmd {
	selected, ok := m.selectedItem()
	if !ok {
		return m.List.NewStatusMessage(styles.ErrorMessage.Render("no item selected"))
	}

	data, err := json.Marshal(selected.Itm)
	if err != nil {
		return func() tea.Msg {
			return ErrorMsg{err}
		}
	}

	writeClipboardText(data)
	return m.List.NewStatusMessage(
		styles.StatusMessage.Render(
			fmt.Sprintf("copied \"%s\" to system clipboard", selected.Itm.Title),
		),
	)
}

// Paste pastes item into current board.
func (m *MenuModel) Paste() tea.Cmd {
	currentBoard, ok := m.selectedBoard()
	if !ok {
		return m.List.NewStatusMessage(styles.ErrorMessage.Render("no board selected"))
	}

	data := readClipboardText()
	var lastItem service.Item
	err := json.Unmarshal(data, &lastItem)
	if err != nil {
		return m.List.NewStatusMessage(
			styles.ErrorMessage.Render(
				fmt.Sprintf("no items in clipboard: %v", err),
			),
		)
	}

	item, err := m.Service.CreateItem(m.ctx, &currentBoard.Board, lastItem.Title, lastItem.Description)
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

// ListItems fetches items in the selected board.
func (m *MenuModel) ListItems() tea.Cmd {
	return func() tea.Msg {
		parentItem, ok := m.selectedBoard()
		if !ok {
			return ErrorMsg{errors.New("no board selected")}
		}

		items, err := m.Service.ListItemsByBoard(m.ctx, &parentItem.Board)
		if err != nil {
			return ErrorMsg{err}
		}
		return ListItemsMsg{
			items,
		}
	}
}

// CreateItem creates a new item.
func (m *MenuModel) CreateItem() tea.Cmd {
	return func() tea.Msg {
		parentItem, ok := m.selectedBoard()
		if !ok {
			return CreateItemMsg{
				Item:  nil,
				Error: errors.New("no board selected"),
			}
		}

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
	m.Copy()
	selected, ok := m.selectedItem()
	if !ok {
		return func() tea.Msg {
			return DeleteItemMsg{Error: errors.New("no item selected")}
		}
	}

	err := m.Service.DeleteItem(m.ctx, &selected.Itm)
	return func() tea.Msg {
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
