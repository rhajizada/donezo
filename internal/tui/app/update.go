package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *AppModel) Push(model tea.Model) {
	// Push a new model onto the stack
	m.ViewStack = append(m.ViewStack, model)
	// Apply the last window size if available
	if m.LastWindowSize != nil {
		m.ApplyWindowSizeToCurrent(*m.LastWindowSize)
	}
}

func (m *AppModel) Pop() {
	// Ensure at least one view remains
	if len(m.ViewStack) > 1 {
		m.ViewStack = m.ViewStack[:len(m.ViewStack)-1]
	}
}

// ApplyWindowSizeToCurrent applies the latest WindowSizeMsg to the current top of the view stack
func (m *AppModel) ApplyWindowSizeToCurrent(msg tea.WindowSizeMsg) tea.Cmd {
	if len(m.ViewStack) == 0 {
		return nil
	}
	// Apply the window size message to the current view
	currentView := m.ViewStack[len(m.ViewStack)-1]
	updatedView, cmd := currentView.Update(msg)
	m.ViewStack[len(m.ViewStack)-1] = updatedView
	return cmd
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		cmd, ok := m.HandleKeyInput(msg)
		if ok {
			return m, cmd
		}
	case tea.WindowSizeMsg:
		// Store the latest window size
		m.LastWindowSize = &msg
		// Apply it to the current view
		return m, m.ApplyWindowSizeToCurrent(msg)

	}

	// Update the active view
	activeView := m.ViewStack[len(m.ViewStack)-1]
	updatedView, cmd := activeView.Update(msg)
	m.ViewStack[len(m.ViewStack)-1] = updatedView

	return m, cmd
}
