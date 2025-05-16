package tags

import (
	"github.com/rhajizada/donezo/internal/tui/styles"
)

func (m MenuModel) View() string {
	return styles.App.Render(m.List.View())
}
