package itemsbytag

import "strings"

// Helper functions for string manipulation.
func truncate(s string, maxLen int, ellipsis string) string {
	if len(s) > maxLen {
		if maxLen-len(ellipsis) > 0 {
			return s[:maxLen-len(ellipsis)] + ellipsis
		}
		return s[:maxLen]
	}
	return s
}

func splitLines(s string) []string {
	return strings.Split(s, "\n")
}
