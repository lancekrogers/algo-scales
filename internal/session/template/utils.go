package template

import "strings"

// sanitizeCommentText sanitizes a text for use in a code comment
func sanitizeCommentText(text string) string {
	// Replace text that might break comments
	text = strings.ReplaceAll(text, "*/", "*\\/")
	return text
}