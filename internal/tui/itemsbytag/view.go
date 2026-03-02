package itemsbytag

import (
	tea "charm.land/bubbletea/v2"

	"github.com/rhajizada/donezo/internal/tui/styles"
)

func (m MenuModel) View() tea.View {
	content := styles.App.Render(m.List.View().Content)
	if m.Context.State != DefaultState {
		content = styles.App.Render(m.Input.View())
	}
	return tea.NewView(content)
}
