package project

import (
	"strings"
)

// sanitizeFileName cleanses the entity names of special characters disallowed by file systems.
func sanitizeFileName(name string) string {
	replacer := strings.NewReplacer(
		"<", "_",
		">", "_",
		"/", "_",
		"\\", "_",
		"?", "_",
		":", "_",
		"*", "_",
		"|", "_",
		"\"", "_")
	return replacer.Replace(name)
}
