package tags

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"

	"github.com/rhajizada/donezo/internal/service"
	"github.com/rhajizada/donezo/internal/tui/styles"
)

//nolint:gochecknoglobals // injectable for tests
var writeClipboardText = func(data []byte) {
	clipboard.Write(clipboard.FmtText, data)
}

func (m *MenuModel) selectedItem() (Item, bool) {
	item, ok := m.List.SelectedItem().(Item)
	return item, ok
}

// ListTags fetches the list of tags from the client.
func (m *MenuModel) ListTags() tea.Cmd {
	return func() tea.Msg {
		data, err := m.Client.ListTags(m.ctx)
		tags := make([]Item, len(data))
		if err != nil {
			return ErrorMsg{err}
		}
		for i, v := range data {
			count, _ := m.Client.CountItemsByTag(m.ctx, v)
			tags[i] = NewItem(v, count)
		}
		return ListTagsMsg{
			tags,
		}
	}
}

// Copy copies tag to system clipboard.
func (m *MenuModel) Copy() tea.Cmd {
	current, ok := m.selectedItem()
	if !ok {
		return m.List.NewStatusMessage(styles.ErrorMessage.Render("no tag selected"))
	}
	currentTag := current.Tag

	items, err := m.Client.ListItemsByTag(m.ctx, currentTag)
	if err != nil {
		return func() tea.Msg {
			return ErrorMsg{err}
		}
	}
	md := service.ItemsToMarkdown(currentTag, *items)
	writeClipboardText([]byte(md))
	return m.List.NewStatusMessage(
		styles.StatusMessage.Render(
			fmt.Sprintf("copied \"%s\" to system clipboard", currentTag),
		),
	)
}

// DeleteTag deletes current selected tag.
func (m *MenuModel) DeleteTag() tea.Cmd {
	return func() tea.Msg {
		selected, ok := m.selectedItem()
		if !ok {
			return DeleteTagMsg{Error: errors.New("no tag selected")}
		}
		err := m.Client.DeleteTag(m.ctx, selected.Tag)
		return DeleteTagMsg{Error: err, Tag: selected.Tag}
	}
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

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

	case ListTagsMsg:
		m.List.SetItems(NewList(msg.Tags))

	case DeleteTagMsg:
		cmd := m.HandleDeleteTag(msg)
		cmds = append(cmds, cmd)
		cmd = m.ListTags()
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
