package app

func (m AppModel) View() string {
	active := m.activeModel()
	if active == nil {
		return ""
	}
	return active.View()
}
