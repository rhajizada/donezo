package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
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
	prevBoard        key.Binding // Added key binding for previous board
	renameItem       key.Binding
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
			key.WithHelp("tab", "next board"),
		),
		prevBoard: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev board"),
		),
		renameItem: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "rename item"),
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

type updateItemMsg struct {
	err error
}

type addItemMsg struct {
	item repository.Item
	err  error
}

// Model defines the UI model.
type Model struct {
	Client       *client.Client
	Boards       []repository.Board
	CurrentBoard int
	List         list.Model
	Keys         *listKeyMap
	DelegateKeys *delegateKeyMap

	// Input states
	renaming     bool
	adding       bool // Added for tracking adding state
	enteringName bool
	enteringDesc bool
	tempName     string
	tempDesc     string
	textInput    textinput.Model
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
			listKeys.insertItem,
			listKeys.deleteItem,
			listKeys.refreshList,
			listKeys.nextBoard,
			listKeys.prevBoard, // Added to help menu
			listKeys.renameItem,
			listKeys.toggleSpinner,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	// Initialize text input
	m.textInput = textinput.New()
	m.textInput.CharLimit = 256
	m.textInput.Width = 50

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
		m.CurrentBoard = len(m.Boards) - 1 // Open the board with the last ID
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

	if m.renaming || m.adding {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				if m.enteringName {
					m.tempName = m.textInput.Value()
					m.enteringName = false
					m.enteringDesc = true
					m.textInput.Placeholder = "Enter item description"
					m.textInput.SetValue("")
					m.textInput.Focus()
				} else if m.enteringDesc {
					m.tempDesc = m.textInput.Value()
					m.enteringDesc = false
					m.textInput.Blur()
					if m.renaming {
						m.renaming = false
						// Proceed to update the item
						index := m.List.Index()
						if index >= 0 && index < len(m.List.Items()) {
							item := m.List.SelectedItem().(Item)
							return m, m.updateItem(&item.Item, m.tempName, m.tempDesc)
						}
					} else if m.adding {
						m.adding = false
						// Proceed to add the item
						board := m.Boards[m.CurrentBoard]
						return m, m.addItem(&board, m.tempName, m.tempDesc)
					}
				}
			case tea.KeyEsc:
				// Cancel the renaming or adding process
				m.renaming = false
				m.adding = false
				m.enteringName = false
				m.enteringDesc = false
				m.textInput.Blur()
			}
		}
		return m, tea.Batch(cmds...)
	}

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

	case updateItemMsg:
		if msg.err != nil {
			cmds = append(cmds, m.List.NewStatusMessage(StatusMessageStyle(fmt.Sprintf("Error updating item: %v", msg.err))))
		} else {
			cmds = append(cmds, m.List.NewStatusMessage(StatusMessageStyle("Item updated")))
			// Refresh the list to show updated item
			return m, m.fetchItems()
		}
		return m, tea.Batch(cmds...)

	case addItemMsg:
		if msg.err != nil {
			cmds = append(cmds, m.List.NewStatusMessage(StatusMessageStyle(fmt.Sprintf("Error adding item: %v", msg.err))))
		} else {
			m.List.InsertItem(0, Item{Item: msg.item})
			cmds = append(cmds, m.List.NewStatusMessage(StatusMessageStyle("Item added")))
		}
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.List.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.Keys.renameItem):
			if len(m.List.Items()) == 0 {
				cmds = append(cmds, m.List.NewStatusMessage(StatusMessageStyle("No item selected")))
				return m, tea.Batch(cmds...)
			}
			m.renaming = true
			m.enteringName = true
			m.textInput.Placeholder = "Enter new item name"
			m.textInput.SetValue("")
			m.textInput.Focus()
			return m, nil

		case key.Matches(msg, m.Keys.insertItem):
			m.adding = true
			m.enteringName = true
			m.textInput.Placeholder = "Enter item name"
			m.textInput.SetValue("")
			m.textInput.Focus()
			return m, nil

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

		case key.Matches(msg, m.Keys.prevBoard):
			if m.CurrentBoard == 0 {
				m.CurrentBoard = len(m.Boards) - 1
			} else {
				m.CurrentBoard--
			}
			return m, m.fetchItems()

		case key.Matches(msg, m.Keys.refreshList):
			return m, m.fetchItems()

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
	if m.renaming || m.adding {
		return appStyle.Render(m.textInput.View())
	}
	return appStyle.Render(m.List.View())
}

// Implement updateItem command
func (m *Model) updateItem(item *repository.Item, newName, newDesc string) tea.Cmd {
	return func() tea.Msg {
		item.Title = newName
		item.Description = newDesc
		_, err := m.Client.UpdateItem(item)
		return updateItemMsg{err: err}
	}
}

// Implement addItem command
func (m *Model) addItem(board *repository.Board, name, desc string) tea.Cmd {
	return func() tea.Msg {
		newItem, err := m.Client.AddItem(board, name, desc)
		if err != nil {
			return addItemMsg{err: err}
		}
		return addItemMsg{item: *newItem}
	}
}
