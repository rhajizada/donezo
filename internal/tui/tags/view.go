package tags

import (
	tea "charm.land/bubbletea/v2"

	"github.com/rhajizada/donezo/internal/tui/styles"
)

func (m MenuModel) View() tea.View {
	return tea.NewView(styles.App.Render(m.List.View()))
}
