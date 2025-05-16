package app

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rhajizada/donezo/internal/tui/boards"
	"github.com/rhajizada/donezo/internal/tui/itemsbyboard"
	"github.com/rhajizada/donezo/internal/tui/itemsbytag"
	"github.com/rhajizada/donezo/internal/tui/tags"
)

func (m *AppModel) NextBoard() (tea.Cmd, bool) {
	var cmd tea.Cmd
	handled := false
	if _, ok := m.ViewStack[len(m.ViewStack)-1].(itemsbyboard.MenuModel); ok {
		// Get the board menu model
		if boardMenu, ok := m.ViewStack[0].(boards.MenuModel); ok {
			index := boardMenu.List.Index()
			boardItems := boardMenu.List.Items()
			if index+1 < len(boardItems) {
				// Move to next board
				boardMenu.List.Select(index + 1)
				// Update the board menu in the view stack
				m.ViewStack[0] = boardMenu

				// Get the new current board
				board := m.GetCurrentBoard()
				if board != nil {
					// Replace the current item menu with a new one for the new board
					itemMenu := itemsbyboard.New(m.ctx, m.Service, board)
					// Replace the top of the view stack
					m.ViewStack[len(m.ViewStack)-1] = itemMenu
					// Apply last window size
					if m.LastWindowSize != nil {
						m.ApplyWindowSizeToCurrent(*m.LastWindowSize)
					}
					cmd = tea.Batch(itemMenu.Init())
					handled = true
				} else {
					// No more boards, create an ErrorMessage
					msg := "no more boards"
					cmd = tea.Batch(
						func() tea.Msg {
							return itemsbyboard.ErrorMsg{
								Error: errors.New(msg),
							}
						},
					)
					handled = true
				}
			} else {
				// No more boards, create an ErrorMessage
				msg := "no more boards"
				cmd = tea.Batch(
					func() tea.Msg {
						return itemsbyboard.ErrorMsg{
							Error: errors.New(msg),
						}
					},
				)
				handled = true
			}
		}
	}
	return cmd, handled
}

func (m *AppModel) PreviousBoard() (tea.Cmd, bool) {
	var cmd tea.Cmd
	handled := false
	// If in item menu, move to previous board
	if _, ok := m.ViewStack[len(m.ViewStack)-1].(itemsbyboard.MenuModel); ok {
		// Get the board menu model
		if boardMenu, ok := m.ViewStack[0].(boards.MenuModel); ok {
			index := boardMenu.List.Index()
			if index-1 >= 0 {
				// Move to previous board
				boardMenu.List.Select(index - 1)
				// Update the board menu in the view stack
				m.ViewStack[0] = boardMenu

				// Get the new current board
				board := m.GetCurrentBoard()
				if board != nil {
					// Replace the current item menu with a new one for the new board
					itemMenu := itemsbyboard.New(m.ctx, m.Service, board)
					// Replace the top of the view stack
					m.ViewStack[len(m.ViewStack)-1] = itemMenu
					// Apply last window size
					if m.LastWindowSize != nil {
						m.ApplyWindowSizeToCurrent(*m.LastWindowSize)
					}
					cmd = tea.Batch(itemMenu.Init())
					handled = true
				}
			} else {
				// No previous boards, create an ErrorMessage
				msg := "no previous boards"
				cmd = tea.Batch(
					func() tea.Msg {
						return itemsbyboard.ErrorMsg{
							Error: errors.New(msg),
						}
					},
				)
				handled = true
			}
		}
	}
	return cmd, handled
}

func (m *AppModel) SelectBoard() (tea.Cmd, bool) {
	var cmd tea.Cmd
	handled := false
	// If on the board view, switch to the item view
	if blist, ok := m.ViewStack[len(m.ViewStack)-1].(boards.MenuModel); ok {
		if blist.State == boards.DefaultState {
			itemMenu := itemsbyboard.New(m.ctx, m.Service, &blist)
			m.Push(itemMenu)
			cmd = tea.Batch(itemMenu.Init())
			handled = true
		}
	}
	return cmd, handled
}

func (m *AppModel) NextTag() (tea.Cmd, bool) {
	var cmd tea.Cmd
	handled := false
	// Ensure the current items view is the tag view.
	if _, ok := m.ViewStack[len(m.ViewStack)-1].(itemsbytag.MenuModel); ok {
		// Get the tag menu model from the main view.
		if tagMenu, ok := m.ViewStack[0].(tags.MenuModel); ok {
			index := tagMenu.List.Index()
			tagItems := tagMenu.List.Items()
			if index+1 < len(tagItems) {
				// Move to next tag in the tags menu.
				tagMenu.List.Select(index + 1)
				m.ViewStack[0] = tagMenu
				// Rebuild the items-by-tag view using the new tag.
				newView := itemsbytag.New(m.ctx, m.Service, &tagMenu)
				m.ViewStack[len(m.ViewStack)-1] = newView
				if m.LastWindowSize != nil {
					m.ApplyWindowSizeToCurrent(*m.LastWindowSize)
				}
				cmd = tea.Batch(newView.Init())
				handled = true
			} else {
				// No more tags.
				msg := "no more tags"
				cmd = tea.Batch(func() tea.Msg {
					return itemsbytag.ErrorMsg{Error: errors.New(msg)}
				})
				handled = true
			}
		}
	}
	return cmd, handled
}

func (m *AppModel) PreviousTag() (tea.Cmd, bool) {
	var cmd tea.Cmd
	handled := false
	if _, ok := m.ViewStack[len(m.ViewStack)-1].(itemsbytag.MenuModel); ok {
		if tagMenu, ok := m.ViewStack[0].(tags.MenuModel); ok {
			index := tagMenu.List.Index()
			if index-1 >= 0 {
				tagMenu.List.Select(index - 1)
				m.ViewStack[0] = tagMenu
				newView := itemsbytag.New(m.ctx, m.Service, &tagMenu)
				m.ViewStack[len(m.ViewStack)-1] = newView
				if m.LastWindowSize != nil {
					m.ApplyWindowSizeToCurrent(*m.LastWindowSize)
				}
				cmd = tea.Batch(newView.Init())
				handled = true
			} else {
				msg := "no previous tags"
				cmd = tea.Batch(func() tea.Msg {
					return itemsbytag.ErrorMsg{Error: errors.New(msg)}
				})
				handled = true
			}
		}
	}
	return cmd, handled
}

func (m *AppModel) SelectTag() (tea.Cmd, bool) {
	var cmd tea.Cmd
	handled := false
	// If the current view is the tags menu...
	if tagMenu, ok := m.ViewStack[len(m.ViewStack)-1].(tags.MenuModel); ok {
		// Create an items-by-tag model.
		// Note: itemsbytag.NewModel expects a pointer to the tag menu.
		itemMenu := itemsbytag.New(m.ctx, m.Service, &tagMenu)
		// Push the new view onto the stack.
		m.Push(itemMenu)
		cmd = tea.Batch(itemMenu.Init())
		handled = true
	}
	return cmd, handled
}

func (m *AppModel) Back() (tea.Cmd, bool) {
	var cmd tea.Cmd
	handled := false
	if len(m.ViewStack) > 0 {
		// Check if the top view is a boards menu.
		if blist, ok := m.ViewStack[len(m.ViewStack)-1].(boards.MenuModel); ok {
			if blist.State == boards.DefaultState {
				m.Pop()
				handled = true
			}
		}
		// Check if the top view is an items-by-board view.
		if ilist, ok := m.ViewStack[len(m.ViewStack)-1].(itemsbyboard.MenuModel); ok {
			if ilist.Context.State == itemsbyboard.DefaultState {
				m.Pop()
				handled = true
			}
		}
		// Also check if the top view is an items-by-tag view.
		if itag, ok := m.ViewStack[len(m.ViewStack)-1].(itemsbytag.MenuModel); ok {
			if itag.Context.State == itemsbytag.DefaultState {
				m.Pop()
				handled = true
			}
		}
	} else {
		cmd = tea.Quit
		handled = true
	}
	return cmd, handled
}

func (m *AppModel) HandleKeyInput(msg tea.KeyMsg) (tea.Cmd, bool) {
	var cmd tea.Cmd
	handled := false

	if len(m.ViewStack) == 1 {
		switch msg.String() {
		case tea.KeyTab.String(), tea.KeyShiftTab.String():
			// Toggle between boards and tags.
			cmd = m.ToggleMainMenu()
			handled = true
		case tea.KeyEnter.String():
			if m.MenuType == MenuBoards {
				cmd, handled = m.SelectBoard()
			} else if m.MenuType == MenuTags {
				cmd, handled = m.SelectTag()
			}
		case tea.KeyBackspace.String():
			cmd, handled = m.Back()
		}
	} else {
		// When an items view is active, dispatch based on its type.
		if _, ok := m.ViewStack[len(m.ViewStack)-1].(itemsbytag.MenuModel); ok {
			switch msg.String() {
			case tea.KeyTab.String():
				cmd, handled = m.NextTag()
			case tea.KeyShiftTab.String():
				cmd, handled = m.PreviousTag()
			case tea.KeyEnter.String():
				// In tag view, you might choose not to do any action on Enter.
				cmd = nil
				handled = false
			case tea.KeyBackspace.String():
				cmd, handled = m.Back()
			}
		} else if _, ok := m.ViewStack[len(m.ViewStack)-1].(itemsbyboard.MenuModel); ok {
			switch msg.String() {
			case tea.KeyTab.String():
				cmd, handled = m.NextBoard()
			case tea.KeyShiftTab.String():
				cmd, handled = m.PreviousBoard()
			case tea.KeyEnter.String():
				cmd, handled = m.SelectBoard()
			case tea.KeyBackspace.String():
				cmd, handled = m.Back()
			}
		}
	}
	return cmd, handled
}

func (m *AppModel) ToggleMainMenu() tea.Cmd {
	var newMain tea.Model
	if m.MenuType == MenuBoards {
		m.MenuType = MenuTags
		newMain = tags.NewModel(m.ctx, m.Service)
	} else {
		m.MenuType = MenuBoards
		newMain = boards.New(m.ctx, m.Service)
	}
	m.ViewStack[0] = newMain
	initCmd := newMain.Init()
	// Force the new model to update its layout if we already have a window size.
	if m.LastWindowSize != nil {
		return tea.Batch(initCmd, m.ApplyWindowSizeToCurrent(*m.LastWindowSize))
	}
	return initCmd
}
