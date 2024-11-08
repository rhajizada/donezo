package items

import "strings"

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
