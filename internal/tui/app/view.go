package app

import tea "charm.land/bubbletea/v2"

func (m AppModel) View() tea.View {
	active := m.activeModel()
	if active == nil {
		return tea.View{}
	}
	view := active.View()
	view.AltScreen = true
	return view
}
