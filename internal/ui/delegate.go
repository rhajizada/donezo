package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Define custom styles
var itemPadding = lipgloss.NewStyle().Padding(0, 2)

// ListDelegate is a fully custom delegate that replicates the default behavior
// but adds a strikethrough to completed items and applies padding.
type ListDelegate struct {
	*list.DefaultDelegate // Embed as a pointer to avoid invalid indirection
	ShowDescription       bool
}

// NewDelegate initializes a new CustomDelegate with default styles.
func NewDelegate() *ListDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true

	return &ListDelegate{
		DefaultDelegate: &delegate,
		ShowDescription: true,
	}
}

// Render overrides the DefaultDelegate's Render method to apply custom styles.
func (d *ListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	itm, ok := item.(Item)
	if !ok {
		return
	}

	title := itm.Title()
	desc := itm.Description()
	completed := itm.Item.Completed

	// Prevent text from exceeding list width
	if m.Width() <= 0 {
		// Short-circuit rendering if width is not set
		return
	}

	textWidth := m.Width() - d.Styles.NormalTitle.GetPaddingLeft() - d.Styles.NormalTitle.GetPaddingRight()
	title = truncate(title, textWidth, "...")
	if d.ShowDescription {
		var lines []string
		for i, line := range splitLines(desc) {
			if i >= m.Height()-2 { // Adjusted to accommodate padding
				break
			}
			lines = append(lines, truncate(line, textWidth, "..."))
		}
		desc = strings.Join(lines, "\n")
	}

	// Determine if the current item is selected
	isSelected := index == m.Index()

	// Apply styles based on selection and completion
	var titleStyle lipgloss.Style
	var descStyle lipgloss.Style

	if isSelected {
		titleStyle = d.Styles.SelectedTitle.Strikethrough(completed)
		descStyle = d.Styles.SelectedDesc.Strikethrough(completed)
	} else {
		titleStyle = d.Styles.NormalTitle.Strikethrough(completed)
		descStyle = d.Styles.NormalDesc.Strikethrough(completed)
	}

	styledTitle := titleStyle.Render(title)
	styledDesc := descStyle.Render(desc)

	// Combine title and description
	var combined string
	if d.ShowDescription {
		combined = fmt.Sprintf("%s\n%s", styledTitle, styledDesc)
	} else {
		combined = styledTitle
	}

	// Apply padding
	combined = itemPadding.Render(combined)

	// Write to the writer
	fmt.Fprint(w, combined) //nolint: errcheck
}

// Override the Update method if necessary (optional)
func (d *ListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	// You can add custom update logic here if needed
	return d.DefaultDelegate.Update(msg, m)
}

// newItemDelegate creates a new CustomDelegate with minimal configuration.
func newItemDelegate(keys *listKeyMap) list.ItemDelegate {
	d := NewDelegate()

	d.ShortHelpFunc = func() []key.Binding {
		return keys.ShortHelp()
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return keys.FullHelp()
	}

	return d
}
