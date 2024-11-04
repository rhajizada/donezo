package ui

import (
	"fmt"
	"math"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/rhajizada/donezo/pkg/client"
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

// listKeyMap defines custom key bindings and implements help.KeyMap interface.
type listKeyMap struct {
	// Custom key bindings
	AddItem          key.Binding
	DeleteItem       key.Binding
	RenameItem       key.Binding
	RefreshList      key.Binding
	NextBoard        key.Binding
	PrevBoard        key.Binding
	ToggleTitleBar   key.Binding
	ToggleStatusBar  key.Binding
	TogglePagination key.Binding
	ToggleComplete   key.Binding
	ToggleHelpMenu   key.Binding
	Quit             key.Binding
}

// newListKeyMap initializes a new listKeyMap with custom bindings.
func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		AddItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		DeleteItem: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete item"),
		),
		RenameItem: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "rename item"),
		),
		RefreshList: key.NewBinding(
			key.WithKeys("R"),
			key.WithHelp("R", "refresh board"),
		),
		NextBoard: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next board"),
		),
		PrevBoard: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev board"),
		),
		ToggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		ToggleStatusBar: key.NewBinding(
			key.WithKeys("B"),
			key.WithHelp("B", "toggle status"),
		),
		TogglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		ToggleComplete: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle complete"),
		),
		// Removed ToggleHelpMenu (H) as per requirement
		ToggleHelpMenu: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q", "quit"),
		),
	}
}

// ShortHelp returns keybindings to be shown in the mini help view.
// It's part of the help.KeyMap interface.
func (k listKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.AddItem,
		k.RefreshList,
		k.NextBoard,
		k.ToggleComplete,
		k.ToggleHelpMenu,
		k.Quit,
	}
}

// FullHelp returns keybindings for the expanded help view.
// It's part of the help.KeyMap interface.
func (k listKeyMap) FullHelp() [][]key.Binding {
	allKeys := []key.Binding{
		k.AddItem,
		k.DeleteItem,
		k.RenameItem,
		k.RefreshList,
		k.NextBoard,
		k.PrevBoard,
		k.ToggleTitleBar,
		k.ToggleStatusBar,
		k.TogglePagination,
		k.ToggleComplete,
		k.ToggleHelpMenu, // Only '?' is included
		k.Quit,
	}

	numCols := 5 // Reduced number of columns to accommodate spacing
	totalKeys := len(allKeys)
	numRows := int(math.Ceil(float64(totalKeys) / float64(numCols)))

	columns := make([][]key.Binding, numCols)
	for i := 0; i < numCols; i++ {
		columns[i] = make([]key.Binding, 0, numRows)
	}

	for i, kb := range allKeys {
		col := i / numRows
		if col >= numCols {
			col = numCols - 1
		}
		columns[col] = append(columns[col], kb)
	}

	// Pad columns with empty key bindings to ensure uniform rows
	for i := 0; i < numCols; i++ {
		for len(columns[i]) < numRows {
			columns[i] = append(columns[i], key.Binding{}) // Empty binding
		}
	}

	// Add spacing between columns
	spacedColumns := make([][]key.Binding, numRows)
	for row := 0; row < numRows; row++ {
		rowBindings := make([]key.Binding, 0, numCols)
		for col := 0; col < numCols; col++ {
			if row < len(columns[col]) {
				rowBindings = append(rowBindings, columns[col][row])
			} else {
				rowBindings = append(rowBindings, key.Binding{})
			}
		}
		spacedColumns[row] = rowBindings
	}

	return spacedColumns
}

// Message types
type errMsg struct {
	err error
}

type boardsLoadedMsg struct{}

type itemsLoadedMsg struct {
	items []client.Item
}

type updateItemMsg struct {
	err error
}

type createItemMsg struct {
	item client.Item
	err  error
}

// Model defines the UI model.
type Model struct {
	Client       *client.Client
	Boards       []client.Board
	CurrentBoard int
	List         list.Model
	Keys         *listKeyMap

	// Input states
	renaming     bool
	adding       bool // Added for tracking adding state
	enteringName bool
	enteringDesc bool
	tempName     string
	tempDesc     string
	textInput    textinput.Model

	// Help
	help     help.Model
	showHelp bool
}

// NewModel initializes a new UI model.
func NewModel(cli *client.Client) *Model {
	listKeys := newListKeyMap()

	// Initialize list with custom delegate
	delegate := newItemDelegate(listKeys)
	l := list.New(nil, delegate, 0, 0)
	l.Title = "Loading boards..."
	l.Styles.Title = titleStyle

	// Initialize text input
	ti := textinput.New()
	ti.CharLimit = 256
	ti.Width = 50

	// Initialize help
	h := help.New()

	m := &Model{
		Client:    cli,
		Keys:      listKeys,
		List:      l,
		textInput: ti,
		help:      h,
		showHelp:  false,
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

// convertToListItems converts repository items to list.Items
func convertToListItems(items []client.Item) []list.Item {
	l := make([]list.Item, len(items))
	for i, item := range items {
		l[i] = Item{Item: item}
	}
	return l
}

// Update handles incoming messages and updates the model accordingly.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle input states: renaming or adding
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

					var placeholder string
					var initialValue string
					if m.renaming {
						index := m.List.Index()
						item := m.List.Items()[index].(Item)
						placeholder = item.Item.Description
						initialValue = item.Item.Description
					} else {
						placeholder = "Enter item description"
						initialValue = ""
					}
					m.textInput.Placeholder = placeholder
					m.textInput.SetValue(initialValue)
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
						return m, m.createItem(&board, m.tempName, m.tempDesc)
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
		m.help.Width = msg.Width / 2 // Adjust help width as needed

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
			// No need to refresh the list here since we're updating the item directly
		}
		return m, tea.Batch(cmds...)

	case createItemMsg:
		if msg.err != nil {
			cmds = append(cmds, m.List.NewStatusMessage(StatusMessageStyle(fmt.Sprintf("Error adding item: %v", msg.err))))
		} else {
			m.List.InsertItem(0, Item{Item: msg.item})
			cmds = append(cmds, m.List.NewStatusMessage(StatusMessageStyle("Item added")))
		}
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		// Handle help menu toggle with only '?'
		if key.Matches(msg, m.Keys.ToggleHelpMenu) {
			m.showHelp = !m.showHelp
			m.help.ShowAll = m.showHelp // Synchronize ShowAll flag
			return m, nil
		}

		// Don't match any of the keys below if we're actively filtering.
		if m.List.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.Keys.RenameItem):
			if len(m.List.Items()) == 0 {
				cmds = append(cmds, m.List.NewStatusMessage(StatusMessageStyle("No item selected")))
				return m, tea.Batch(cmds...)
			}
			m.renaming = true
			m.enteringName = true
			index := m.List.Index()
			item := m.List.Items()[index].(Item)

			m.textInput.Placeholder = item.Item.Title
			m.textInput.SetValue(item.Item.Title)
			m.textInput.Focus()
			return m, nil

		case key.Matches(msg, m.Keys.AddItem):
			m.adding = true
			m.enteringName = true
			m.textInput.Placeholder = "Enter item name"
			m.textInput.SetValue("")
			m.textInput.Focus()
			return m, nil

		case key.Matches(msg, m.Keys.RefreshList):
			return m, m.fetchItems()

		case key.Matches(msg, m.Keys.ToggleTitleBar):
			v := !m.List.ShowTitle()
			m.List.SetShowTitle(v)
			m.List.SetShowFilter(v)
			m.List.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.Keys.ToggleStatusBar):
			m.List.SetShowStatusBar(!m.List.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.Keys.TogglePagination):
			m.List.SetShowPagination(!m.List.ShowPagination())
			return m, nil

		case key.Matches(msg, m.Keys.NextBoard):
			if len(m.Boards) == 0 {
				cmds = append(cmds, m.List.NewStatusMessage(StatusMessageStyle("No boards available")))
				return m, tea.Batch(cmds...)
			}
			m.CurrentBoard = (m.CurrentBoard + 1) % len(m.Boards)
			return m, m.fetchItems()

		case key.Matches(msg, m.Keys.PrevBoard):
			if len(m.Boards) == 0 {
				cmds = append(cmds, m.List.NewStatusMessage(StatusMessageStyle("No boards available")))
				return m, tea.Batch(cmds...)
			}
			if m.CurrentBoard == 0 {
				m.CurrentBoard = len(m.Boards) - 1
			} else {
				m.CurrentBoard--
			}
			return m, m.fetchItems()

		case key.Matches(msg, m.Keys.DeleteItem):
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

		case key.Matches(msg, m.Keys.ToggleComplete):
			index := m.List.Index()
			if index >= 0 && index < len(m.List.Items()) {
				item := m.List.Items()[index].(Item)
				return m, m.toggleItemCompletion(index, item)
			}
		}
	}

	// Update the list
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the UI.
func (m *Model) View() string {
	if m.renaming || m.adding {
		return appStyle.Render(m.textInput.View())
	}

	view := m.List.View()

	if m.showHelp {
		helpView := m.help.View(m.Keys)
		view += "\n\n" + appStyle.Render(helpView)
	}

	return appStyle.Render(view)
}

// Implement updateItem command
func (m *Model) updateItem(item *client.Item, newName, newDesc string) tea.Cmd {
	return func() tea.Msg {
		item.Title = newName
		item.Description = newDesc
		_, err := m.Client.UpdateItem(item)
		return updateItemMsg{err: err}
	}
}

// Implement createItem command
func (m *Model) createItem(board *client.Board, name, desc string) tea.Cmd {
	return func() tea.Msg {
		newItem, err := m.Client.CreateItem(board, name, desc)
		if err != nil {
			return createItemMsg{err: err}
		}
		return createItemMsg{item: *newItem}
	}
}

// Implement toggleItemCompletion command
func (m *Model) toggleItemCompletion(index int, item Item) tea.Cmd {
	// Toggle the completion status
	item.Completed = !item.Completed
	// Update the item in the list
	m.List.SetItem(index, item)
	return func() tea.Msg {
		_, err := m.Client.UpdateItem(&item.Item)
		if err != nil {
			return errMsg{err}
		}
		return updateItemMsg{err: nil}
	}
}
