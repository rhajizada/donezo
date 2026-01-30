package itemsbyboard

import (
	"fmt"
	"io"
	"strings"

	"github.com/rhajizada/donezo/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"

	"github.com/rhajizada/donezo/internal/tui/itemlist"
)

// ListDelegate is a fully custom delegate that replicates the default behavior
// but adds a strikethrough to completed items and applies padding.
type ListDelegate struct {
	*itemlist.DefaultDelegate // Embed as a pointer to avoid invalid indirection
}

// NewDelegate initializes a new CustomDelegate with default styles.
func NewDelegate() *ListDelegate {
	delegate := itemlist.NewDefaultDelegate()

	return &ListDelegate{
		DefaultDelegate: &delegate,
	}
}

// Render overrides the DefaultDelegate's Render method to apply custom styles.
func (d *ListDelegate) Render(w io.Writer, m itemlist.Model, index int, item itemlist.Item) {
	selected, ok := item.(Item)
	if !ok {
		return
	}

	title := selected.Itm.Title
	desc := selected.Itm.Description
	completed := selected.Itm.Completed

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
			if i >= m.Height()-2 {
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
	combined = styles.Item.Render(combined)

	// Write to the writer
	fmt.Fprint(w, combined)
}
