package styles

import "github.com/charmbracelet/lipgloss"

const (
	appMarginVertical   = 1
	appMarginHorizontal = 2
	footerMargin        = 2
)

//nolint:gochecknoglobals // shared lipgloss styles reused across TUI views
var (
	App = lipgloss.NewStyle().
		Margin(appMarginVertical, appMarginHorizontal)

	StatusMessage = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"})

	ErrorMessage = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#FB4A8A", Dark: "#FB4A8A"})

	Item = lipgloss.NewStyle().
		Padding(0, 0)

	Footer = lipgloss.NewStyle().
		Margin(0, footerMargin).
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4"))
)
