package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/rhajizada/donezo/client"
	"github.com/rhajizada/donezo/internal/repository"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)
)

// StatusMessageStyle styles status messages in the UI.
func StatusMessageStyle(msg string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
		Render(msg)
}

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
	deleteItem       key.Binding
	refreshList      key.Binding
	nextBoard        key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		deleteItem: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete item"),
		),
		refreshList: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "refresh list"),
		),
		nextBoard: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "switch board"),
		),
		toggleSpinner: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("B"),
			key.WithHelp("B", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

// Message types
type errMsg struct {
	err error
}

type boardsLoadedMsg struct{}

type itemsLoadedMsg struct {
	items []repository.Item
}

// Model defines the UI model.
type Model struct {
	Client       *client.Client
	Boards       []repository.Board
	CurrentBoard int
	List         list.Model
	Keys         *listKeyMap
	DelegateKeys *delegateKeyMap
}

// NewModel initializes a new UI model.
func NewModel(cli *client.Client) *Model {
	var (
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
	)

	m := &Model{
		Client:       cli,
		Keys:         listKeys,
		DelegateKeys: delegateKeys,
	}

	// Initialize list
	delegate := newItemDelegate(delegateKeys)
	m.List = list.New(nil, delegate, 0, 0)
	m.List.Title = "Loading boards..."
	m.List.Styles.Title = titleStyle
	m.List.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.insertItem,
			listKeys.deleteItem,
			listKeys.refreshList,
			listKeys.nextBoard,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	return m
}

// Init initializes the UI.
func (m *Model) Init() tea.Cmd {
	return m.fetchBoards()
}

// fetchBoards fetches the list of boards from the client.
func (m *Model) fetchBoards() tea.Cmd {
	return func() tea.Msg {
		boards, err := m.Client.ListBoards()
		if err != nil {
			return errMsg{err}
		}
		m.Boards = *boards
		if len(m.Boards) == 0 {
			return errMsg{fmt.Errorf("no boards available")}
		}
		m.CurrentBoard = 0
		return boardsLoadedMsg{}
	}
}

// fetchItems fetches items for the current board.
func (m *Model) fetchItems() tea.Cmd {
	board := m.Boards[m.CurrentBoard]
	m.List.Title = board.Name
	return func() tea.Msg {
		items, err := m.Client.ListItems(&board)
		if err != nil {
			return errMsg{err}
		}
		return itemsLoadedMsg{items: *items}
	}
}

func convertToListItems(items []repository.Item) []list.Item {
	l := make([]list.Item, len(items))
	for i, item := range items {
		l[i] = Item{Item: item}
	}
	return l
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)

	case errMsg:
		cmds = append(cmds, m.List.NewStatusMessage(StatusMessageStyle(fmt.Sprintf("Error: %v", msg.err))))
		return m, tea.Batch(cmds...)

	case boardsLoadedMsg:
		return m, m.fetchItems()

	case itemsLoadedMsg:
		m.List.SetItems(convertToListItems(msg.items))
		// Optionally, you can add a status message or perform additional actions

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.List.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.Keys.toggleSpinner):
			cmd := m.List.ToggleSpinner()
			return m, cmd

		case key.Matches(msg, m.Keys.toggleTitleBar):
			v := !m.List.ShowTitle()
			m.List.SetShowTitle(v)
			m.List.SetShowFilter(v)
			m.List.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.Keys.toggleStatusBar):
			m.List.SetShowStatusBar(!m.List.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.Keys.togglePagination):
			m.List.SetShowPagination(!m.List.ShowPagination())
			return m, nil

		case key.Matches(msg, m.Keys.toggleHelpMenu):
			m.List.SetShowHelp(!m.List.ShowHelp())
			return m, nil

		case key.Matches(msg, m.Keys.nextBoard):
			m.CurrentBoard = (m.CurrentBoard + 1) % len(m.Boards)
			return m, m.fetchItems()

		case key.Matches(msg, m.Keys.refreshList):
			return m, m.fetchItems()

		case key.Matches(msg, m.Keys.insertItem):
			board := m.Boards[m.CurrentBoard]
			// For simplicity, we use static values; you can modify to accept user input.
			newItemTitle := "New Item"
			newItemDescription := "Description"
			newItem, err := m.Client.AddItem(&board, newItemTitle, newItemDescription)
			if err != nil {
				cmds = append(cmds, m.List.NewStatusMessage(StatusMessageStyle(fmt.Sprintf("Error adding item: %v", err))))
			} else {
				m.List.InsertItem(0, Item{Item: *newItem})
				cmds = append(cmds, m.List.NewStatusMessage(StatusMessageStyle("Item added")))
			}
			return m, tea.Batch(cmds...)

		case key.Matches(msg, m.Keys.deleteItem):
			index := m.List.Index()
			if index >= 0 && index < len(m.List.Items()) {
				item := m.List.SelectedItem().(Item)
				err := m.Client.DeleteItem(&item.Item)
				if err != nil {
					cmds = append(cmds, m.List.NewStatusMessage(StatusMessageStyle(fmt.Sprintf("Error deleting item: %v", err))))
				} else {
					m.List.RemoveItem(index)
					cmds = append(cmds, m.List.NewStatusMessage(StatusMessageStyle("Item deleted")))
				}
			}
			return m, tea.Batch(cmds...)
		}
	}

	// Update the list
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return appStyle.Render(m.List.View())
}
