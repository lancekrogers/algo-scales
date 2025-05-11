package view

import (
	"bytes"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

// SyntaxHighlighter provides code syntax highlighting functionality
type SyntaxHighlighter struct {
	// Default style to use
	defaultStyle string
}

// NewSyntaxHighlighter creates a new syntax highlighter
func NewSyntaxHighlighter(style string) *SyntaxHighlighter {
	return &SyntaxHighlighter{
		defaultStyle: style,
	}
}

// Highlight returns syntax highlighted code for the terminal
func (h *SyntaxHighlighter) Highlight(code, language string) (string, error) {
	// Get the lexer for the language
	l := lexers.Get(language)
	if l == nil {
		// Try to match by extension
		l = lexers.Match(language)
		if l == nil {
			// Fallback to plain text
			l = lexers.Get("text")
		}
	}
	l = chroma.Coalesce(l)

	// Get the style
	s := styles.Get(h.defaultStyle)
	if s == nil {
		s = styles.Get("monokai")
	}

	// Get the formatter
	f := formatters.Get("terminal")
	if f == nil {
		f = formatters.Get("terminal16m")
	}

	// Create an output buffer
	var buf bytes.Buffer

	// Tokenize and format the code
	it, err := l.Tokenise(nil, code)
	if err != nil {
		return "", err
	}

	err = f.Format(&buf, s, it)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// GetLanguageExtension returns the file extension for a given language
func GetLanguageExtension(language string) string {
	switch strings.ToLower(language) {
	case "go", "golang":
		return "go"
	case "python", "py":
		return "py"
	case "javascript", "js":
		return "js"
	case "typescript", "ts":
		return "ts"
	case "java":
		return "java"
	case "c++", "cpp":
		return "cpp"
	case "c":
		return "c"
	case "c#", "csharp":
		return "cs"
	case "php":
		return "php"
	case "ruby", "rb":
		return "rb"
	case "rust", "rs":
		return "rs"
	case "kotlin", "kt":
		return "kt"
	case "swift":
		return "swift"
	default:
		return "txt"
	}
}

// GetSupportedLanguages returns a list of supported programming languages
func GetSupportedLanguages() []string {
	return []string{
		"go",
		"python",
		"javascript",
		"typescript",
		"java",
		"c++",
		"c",
		"c#",
		"php",
		"ruby",
		"rust",
		"kotlin",
		"swift",
	}
}

// GetLanguageDisplayName returns a user-friendly display name for a language
func GetLanguageDisplayName(language string) string {
	switch strings.ToLower(language) {
	case "go", "golang":
		return "Go"
	case "python", "py":
		return "Python"
	case "javascript", "js":
		return "JavaScript"
	case "typescript", "ts":
		return "TypeScript"
	case "java":
		return "Java"
	case "c++", "cpp":
		return "C++"
	case "c":
		return "C"
	case "c#", "csharp", "cs":
		return "C#"
	case "php":
		return "PHP"
	case "ruby", "rb":
		return "Ruby"
	case "rust", "rs":
		return "Rust"
	case "kotlin", "kt":
		return "Kotlin"
	case "swift":
		return "Swift"
	default:
		return strings.Title(language)
	}
}