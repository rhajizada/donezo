package itemsbytag

import (
	"github.com/rhajizada/donezo/internal/tui/styles"
)

func (m MenuModel) View() string {
	if m.Context.State != DefaultState {
		return styles.App.Render(m.Input.View())
	}
	return styles.App.Render(m.List.View())
}
