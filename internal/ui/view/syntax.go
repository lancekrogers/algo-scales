package view

import (
	"github.com/lancekrogers/algo-scales/internal/common/highlight"
)

// Re-export the syntax highlighter from the common package
var (
	// SyntaxHighlighter is the type from the common package
	SyntaxHighlighter = highlight.SyntaxHighlighter

	// NewSyntaxHighlighter creates a new syntax highlighter
	NewSyntaxHighlighter = highlight.NewSyntaxHighlighter

	// GetLanguageExtension returns the file extension for a language
	GetLanguageExtension = highlight.GetLanguageExtension

	// GetSupportedLanguages returns a list of supported languages
	GetSupportedLanguages = highlight.GetSupportedLanguages

	// GetLanguageDisplayName returns a user-friendly display name
	GetLanguageDisplayName = highlight.GetLanguageDisplayName
)