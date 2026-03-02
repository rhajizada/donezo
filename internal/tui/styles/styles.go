package styles

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
)

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
			Foreground(compat.AdaptiveColor{Light: lipgloss.Color("#04B575"), Dark: lipgloss.Color("#04B575")})

	ErrorMessage = lipgloss.NewStyle().
			Foreground(compat.AdaptiveColor{Light: lipgloss.Color("#FB4A8A"), Dark: lipgloss.Color("#FB4A8A")})

	Item = lipgloss.NewStyle().
		Padding(0, 0)

	Footer = lipgloss.NewStyle().
		Margin(0, footerMargin).
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4"))
)
