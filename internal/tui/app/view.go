package app

func (m AppModel) View() string {
	// Render the active view
	return m.ViewStack[len(m.ViewStack)-1].View()
}
