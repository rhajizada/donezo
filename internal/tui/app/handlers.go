package app

import (
	"fmt"

	"github.com/rhajizada/donezo/internal/tui/boards"
	"github.com/rhajizada/donezo/internal/tui/items"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *AppModel) NextBoard() (tea.Cmd, bool) {
	var cmd tea.Cmd
	handled := false
	if _, ok := m.ViewStack[len(m.ViewStack)-1].(items.ItemMenuModel); ok {
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
					itemMenu := items.NewModel(m.Client, board)
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
							return items.ErrorMsg{
								Error: fmt.Errorf(msg),
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
						return items.ErrorMsg{
							Error: fmt.Errorf(msg),
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
	if _, ok := m.ViewStack[len(m.ViewStack)-1].(items.ItemMenuModel); ok {
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
					itemMenu := items.NewModel(m.Client, board)
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
						return items.ErrorMsg{
							Error: fmt.Errorf(msg),
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
		board := m.GetCurrentBoard()
		if board != nil && blist.State == boards.DefaultState {
			itemMenu := items.NewModel(m.Client, board)
			m.Push(itemMenu)
			cmd = tea.Batch(itemMenu.Init())
			handled = true
		}
	}
	return cmd, handled
}

func (m *AppModel) Back() (tea.Cmd, bool) {
	var cmd tea.Cmd
	handled := false
	if len(m.ViewStack) > 1 {
		{
			m.Pop()
		}
		handled = true
	} else {
		cmd = tea.Quit
		handled = true
	}
	return cmd, handled
}

func (m *AppModel) HandleKeyInput(msg tea.KeyMsg) (tea.Cmd, bool) {
	var cmd tea.Cmd
	handled := false
	switch msg.String() {
	case tea.KeyTab.String():
		cmd, handled = m.NextBoard()
	case tea.KeyShiftTab.String():
		cmd, handled = m.PreviousBoard()
	case tea.KeyEnter.String():
		cmd, handled = m.SelectBoard()
	case "b":
		cmd, handled = m.Back()
	}
	return cmd, handled
}
