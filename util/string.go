package util

import "strings"

// After returns the substring after the first instance of the key
func After(value string, key string) string {
	// Get substring after a string.
	pos := strings.Index(value, key)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(key)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:]
}
