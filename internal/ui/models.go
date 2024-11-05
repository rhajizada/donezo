package ui

import (
	"fmt"

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

// InputState represents the current input state of the UI.
type InputState int

const (
	StateIdle InputState = iota
	StateAddingName
	StateAddingDesc
	StateRenamingName
	StateRenamingDesc
)

// Model defines the UI model.
type Model struct {
	Client       *client.Client
	Boards       []client.Board
	CurrentBoard int
	List         list.Model
	Keys         *Keymap

	// Input states
	inputState  InputState
	tempName    string
	tempDesc    string
	textInput   textinput.Model
	currentItem *Item

	// Help
	help     help.Model
	showHelp bool
}

// NewModel initializes a new UI model.
func NewModel(cli *client.Client) *Model {
	listKeys := newKeymap()

	// Initialize list with custom delegate
	delegate := NewDelegate(listKeys)
	l := list.New(nil, delegate, 0, 0)
	l.Title = "Loading boards..."
	l.Styles.Title = titleStyle

	// Initialize text input
	ti := textinput.New()
	ti.CharLimit = 256
	ti.Width = 50

	m := &Model{
		Client:      cli,
		Keys:        listKeys,
		List:        l,
		textInput:   ti,
		showHelp:    false,
		inputState:  StateIdle,
		currentItem: nil,
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

	// Handle input states
	if m.isInputState() {
		m.textInput, cmds = m.handleInputState(msg)
		return m, tea.Batch(cmds...)
	}

	// Handle general messages
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cmd := m.handleWindowSize(msg)
		cmds = append(cmds, cmd)

	case errMsg:
		cmd := m.handleError(msg)
		cmds = append(cmds, cmd)

	case boardsLoadedMsg:
		cmd := m.fetchItems()
		cmds = append(cmds, cmd)

	case itemsLoadedMsg:
		m.List.SetItems(convertToListItems(msg.items))

	case updateItemMsg:
		cmd := m.handleUpdateItem(msg)
		cmds = append(cmds, cmd)

	case createItemMsg:
		cmd := m.handleCreateItem(msg)
		cmds = append(cmds, cmd)

	case tea.KeyMsg:
		cmd := m.handleGeneralKey(msg)
		cmds = append(cmds, cmd)
	}

	// Update the list
	listModel, listCmd := m.List.Update(msg)
	m.List = listModel
	cmds = append(cmds, listCmd)

	return m, tea.Batch(cmds...)
}

// View renders the UI.
func (m *Model) View() string {
	if m.isInputState() {
		return appStyle.Render(m.textInput.View())
	}

	view := m.List.View()

	return appStyle.Render(view)
}

// isInputState checks if the current state is an input state.
func (m *Model) isInputState() bool {
	return m.inputState == StateAddingName || m.inputState == StateAddingDesc ||
		m.inputState == StateRenamingName || m.inputState == StateRenamingDesc
}

// handleInputState processes messages when in an input state.
func (m *Model) handleInputState(msg tea.Msg) (textinput.Model, []tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	// Only handle key messages in input states
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyEnter:
			switch m.inputState {
			case StateAddingName:
				m.tempName = m.textInput.Value()
				m.inputState = StateAddingDesc
				m.textInput.Placeholder = "Enter item description"
				m.textInput.SetValue("")
				m.textInput.Focus()
			case StateAddingDesc:
				m.tempDesc = m.textInput.Value()
				m.inputState = StateIdle
				m.textInput.Blur()
				cmds = append(cmds, m.createItem())
			case StateRenamingName:
				m.tempName = m.textInput.Value()
				m.inputState = StateRenamingDesc
				if m.currentItem != nil {
					m.textInput.Placeholder = m.currentItem.Item.Description
					m.textInput.SetValue(m.currentItem.Item.Description)
					m.textInput.Focus()
				} else {
					m.textInput.Placeholder = "Enter item description"
					m.textInput.SetValue("")
					m.textInput.Focus()
				}
			case StateRenamingDesc:
				m.tempDesc = m.textInput.Value()
				m.inputState = StateIdle
				m.textInput.Blur()
				cmds = append(cmds, m.updateItem())
			}

		case tea.KeyEsc:
			// Cancel the current operation
			m.inputState = StateIdle
			m.textInput.Blur()
			m.currentItem = nil
		}
	}

	return m.textInput, cmds
}

// handleWindowSize processes window size messages.
func (m *Model) handleWindowSize(msg tea.WindowSizeMsg) tea.Cmd {
	h, v := appStyle.GetFrameSize()

	// Define fixed help height
	helpHeight := 10 // Adjust based on expected help content

	if m.showHelp {
		m.List.SetSize(msg.Width-h, msg.Height-v-helpHeight)
	} else {
		m.List.SetSize(msg.Width-h, msg.Height-v)
	}
	m.help.Width = msg.Width / 2 // Adjust help width as needed

	return nil
}

// handleError processes error messages.
func (m *Model) handleError(msg errMsg) tea.Cmd {
	return m.List.NewStatusMessage(StatusMessageStyle(fmt.Sprintf("Error: %v", msg.err)))
}

// handleUpdateItem processes item update messages.
func (m *Model) handleUpdateItem(msg updateItemMsg) tea.Cmd {
	if msg.err != nil {
		return m.List.NewStatusMessage(StatusMessageStyle(fmt.Sprintf("Error updating item: %v", msg.err)))
	}

	if m.currentItem != nil {
		// Update the item in the list
		index := m.findItemIndex(m.currentItem)
		if index >= 0 {
			m.List.SetItem(index, *m.currentItem)
		}
		m.List.NewStatusMessage(StatusMessageStyle("Item updated"))
		m.currentItem = nil
	} else {
		m.List.NewStatusMessage(StatusMessageStyle("Item updated"))
	}

	return nil
}

// handleCreateItem processes item creation messages.
func (m *Model) handleCreateItem(msg createItemMsg) tea.Cmd {
	if msg.err != nil {
		return m.List.NewStatusMessage(StatusMessageStyle(fmt.Sprintf("Error adding item: %v", msg.err)))
	}

	m.List.InsertItem(0, Item{Item: msg.item})
	return m.List.NewStatusMessage(StatusMessageStyle("Item added"))
}

// handleGeneralKey processes general key messages.
func (m *Model) handleGeneralKey(msg tea.KeyMsg) tea.Cmd {
	// Handle help menu toggle with only '?'
	if key.Matches(msg, m.Keys.ToggleHelpMenu) {
		m.showHelp = !m.showHelp
		m.help.ShowAll = m.showHelp // Synchronize ShowAll flag
		return nil
	}

	// Don't match any of the keys below if we're actively filtering.
	if m.List.FilterState() == list.Filtering {
		return nil
	}

	switch {
	case key.Matches(msg, m.Keys.RenameItem):
		return m.initiateRename()

	case key.Matches(msg, m.Keys.AddItem):
		return m.initiateAdd()

	case key.Matches(msg, m.Keys.RefreshList):
		return m.fetchItems()

	case key.Matches(msg, m.Keys.ToggleTitleBar):
		m.toggleTitleBar()
		return nil

	case key.Matches(msg, m.Keys.ToggleStatusBar):
		m.List.SetShowStatusBar(!m.List.ShowStatusBar())
		return nil

	case key.Matches(msg, m.Keys.TogglePagination):
		m.List.SetShowPagination(!m.List.ShowPagination())
		return nil

	case key.Matches(msg, m.Keys.NextBoard):
		return m.switchBoard(1)

	case key.Matches(msg, m.Keys.PrevBoard):
		return m.switchBoard(-1)

	case key.Matches(msg, m.Keys.DeleteItem):
		return m.deleteItem()

	case key.Matches(msg, m.Keys.ToggleComplete):
		return m.toggleCompletion()
	}

	return nil
}

// initiateRename starts the renaming process for the selected item.
func (m *Model) initiateRename() tea.Cmd {
	if len(m.List.Items()) == 0 {
		return m.List.NewStatusMessage(StatusMessageStyle("No item selected"))
	}

	m.inputState = StateRenamingName
	index := m.List.Index()
	item := m.List.Items()[index].(Item)
	m.currentItem = &item

	m.textInput.Placeholder = item.Item.Title
	m.textInput.SetValue(item.Item.Title)
	m.textInput.Focus()

	return nil
}

// initiateAdd starts the adding process for a new item.
func (m *Model) initiateAdd() tea.Cmd {
	m.inputState = StateAddingName
	m.textInput.Placeholder = "Enter item name"
	m.textInput.SetValue("")
	m.textInput.Focus()
	return nil
}

// toggleTitleBar toggles the visibility of the title bar.
func (m *Model) toggleTitleBar() {
	v := !m.List.ShowTitle()
	m.List.SetShowTitle(v)
	m.List.SetShowFilter(v)
	m.List.SetFilteringEnabled(v)
}

// switchBoard changes the current board by an offset (+1 for next, -1 for previous).
func (m *Model) switchBoard(offset int) tea.Cmd {
	if len(m.Boards) == 0 {
		return m.List.NewStatusMessage(StatusMessageStyle("No boards available"))
	}

	m.CurrentBoard = (m.CurrentBoard + offset + len(m.Boards)) % len(m.Boards)
	return m.fetchItems()
}

// deleteItem deletes the selected item.
func (m *Model) deleteItem() tea.Cmd {
	index := m.List.Index()
	if index < 0 || index >= len(m.List.Items()) {
		return m.List.NewStatusMessage(StatusMessageStyle("No item selected"))
	}

	item := m.List.SelectedItem().(Item)
	err := m.Client.DeleteItem(&item.Item)
	if err != nil {
		return m.List.NewStatusMessage(StatusMessageStyle(fmt.Sprintf("Error deleting item: %v", err)))
	}

	m.List.RemoveItem(index)
	return m.List.NewStatusMessage(StatusMessageStyle("Item deleted"))
}

// toggleCompletion toggles the completion status of the selected item.
func (m *Model) toggleCompletion() tea.Cmd {
	index := m.List.Index()
	if index < 0 || index >= len(m.List.Items()) {
		return m.List.NewStatusMessage(StatusMessageStyle("No item selected"))
	}

	item := m.List.Items()[index].(Item)
	item.Completed = !item.Completed
	m.List.SetItem(index, item)

	return func() tea.Msg {
		_, err := m.Client.UpdateItem(&item.Item)
		if err != nil {
			return errMsg{err}
		}
		return updateItemMsg{err: nil}
	}
}

// Implement updateItem command
func (m *Model) updateItem() tea.Cmd {
	if m.currentItem == nil {
		return func() tea.Msg {
			return errMsg{fmt.Errorf("no item selected for update")}
		}
	}

	item := m.currentItem
	item.Item.Title = m.tempName
	item.Item.Description = m.tempDesc

	return func() tea.Msg {
		_, err := m.Client.UpdateItem(&item.Item)
		return updateItemMsg{err: err}
	}
}

// Implement createItem command
func (m *Model) createItem() tea.Cmd {
	board := m.Boards[m.CurrentBoard]
	return func() tea.Msg {
		newItem, err := m.Client.CreateItem(&board, m.tempName, m.tempDesc)
		if err != nil {
			return createItemMsg{err: err}
		}
		return createItemMsg{item: *newItem}
	}
}

// findItemIndex finds the index of the currentItem in the list.
func (m *Model) findItemIndex(target *Item) int {
	for i, listItem := range m.List.Items() {
		item, ok := listItem.(Item)
		if !ok {
			continue
		}
		if item.ID == target.ID {
			return i
		}
	}
	return -1
}
