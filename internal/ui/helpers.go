package ui

import (
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/key"
)

// Helper functions for string manipulation
func truncate(s string, max int, ellipsis string) string {
	if len(s) > max {
		if max-len(ellipsis) > 0 {
			return s[:max-len(ellipsis)] + ellipsis
		}
		return s[:max]
	}
	return s
}

func splitLines(s string) []string {
	return strings.Split(s, "\n")
}

func joinLines(lines []string) string {
	return strings.Join(lines, "\n")
}

func formatKeys(keys []key.Binding, numRows int) [][]key.Binding {
	totalKeys := len(keys)
	numCols := int(math.Ceil(float64(totalKeys) / float64(numRows)))

	rows := make([][]key.Binding, numRows)
	for i := 0; i < numRows; i++ {
		start := i * numCols
		end := start + numCols
		if end > totalKeys {
			end = totalKeys
		}
		rows[i] = keys[start:end]
		// Pad the row with empty bindings if necessary
		for i := len(rows[i]) - 1; i < numCols+1; i++ {
			rows[i] = append(rows[i], key.Binding{}) // Empty binding
		}
	}

	return rows
}
