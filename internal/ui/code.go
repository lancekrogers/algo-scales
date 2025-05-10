// Syntax highlighting for code display

package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// SyntaxHighlighter provides simple syntax highlighting for code
type SyntaxHighlighter struct {
	// Color styles for different syntax elements
	keywordStyle    lipgloss.Style
	stringStyle     lipgloss.Style
	commentStyle    lipgloss.Style
	numberStyle     lipgloss.Style
	functionStyle   lipgloss.Style
	variableStyle   lipgloss.Style
	operatorStyle   lipgloss.Style
	typeStyle       lipgloss.Style
	defaultStyle    lipgloss.Style
	backgroundColor string
}

// NewSyntaxHighlighter creates a new syntax highlighter
func NewSyntaxHighlighter() *SyntaxHighlighter {
	bg := "#1E1E1E" // Dark background

	return &SyntaxHighlighter{
		keywordStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("#569CD6")), // Blue
		stringStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#CE9178")), // Orange
		commentStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("#6A9955")), // Green
		numberStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#B5CEA8")), // Light green
		functionStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("#DCDCAA")), // Yellow
		variableStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("#9CDCFE")), // Light blue
		operatorStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("#D4D4D4")), // White
		typeStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("#4EC9B0")), // Teal
		defaultStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("#D4D4D4")), // White
		backgroundColor: bg,
	}
}

// HighlightGo highlights Go code
func (h *SyntaxHighlighter) HighlightGo(code string) string {
	lines := strings.Split(code, "\n")
	highlightedLines := make([]string, len(lines))

	// Keywords in Go
	keywords := map[string]bool{
		"func":      true,
		"package":   true,
		"import":    true,
		"type":      true,
		"struct":    true,
		"interface": true,
		"map":       true,
		"chan":      true,
		"const":     true,
		"var":       true,
		"if":        true,
		"else":      true,
		"for":       true,
		"range":     true,
		"switch":    true,
		"case":      true,
		"default":   true,
		"break":     true,
		"continue":  true,
		"return":    true,
		"go":        true,
		"defer":     true,
		"select":    true,
		"make":      true,
		"new":       true,
		"true":      true,
		"false":     true,
		"nil":       true,
	}

	// Types in Go
	types := map[string]bool{
		"string":  true,
		"int":     true,
		"int8":    true,
		"int16":   true,
		"int32":   true,
		"int64":   true,
		"uint":    true,
		"uint8":   true,
		"uint16":  true,
		"uint32":  true,
		"uint64":  true,
		"uintptr": true,
		"byte":    true,
		"rune":    true,
		"float32": true,
		"float64": true,
		"bool":    true,
		"error":   true,
	}

	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			highlightedLines[i] = line
			continue
		}

		// Handle comments
		if commentIdx := strings.Index(line, "//"); commentIdx != -1 {
			beforeComment := line[:commentIdx]
			comment := line[commentIdx:]
			highlightedLines[i] = h.highlightGoParts(beforeComment, keywords, types) + h.commentStyle.Render(comment)
			continue
		}

		highlightedLines[i] = h.highlightGoParts(line, keywords, types)
	}

	return strings.Join(highlightedLines, "\n")
}

// highlightGoParts highlights parts of a Go code line
func (h *SyntaxHighlighter) highlightGoParts(line string, keywords, types map[string]bool) string {
	var result strings.Builder
	var token strings.Builder
	var inString bool
	var inRune bool
	var stringQuote rune

	for _, r := range line {
		// Handle strings
		if (r == '"' || r == '\'') && !inRune && !inString {
			// If we have a token, highlight it before starting the string
			if token.Len() > 0 {
				result.WriteString(h.highlightGoToken(token.String(), keywords, types))
				token.Reset()
			}
			inString = true
			stringQuote = r
			token.WriteRune(r)
			continue
		}

		if inString && r == stringQuote {
			// End of string
			token.WriteRune(r)
			result.WriteString(h.stringStyle.Render(token.String()))
			token.Reset()
			inString = false
			continue
		}

		if inString {
			token.WriteRune(r)
			continue
		}

		// Handle operators and delimiters
		if strings.ContainsRune("+-*/=<>!&|^~();:,.[]{}%", r) {
			// If we have a token, highlight it before the operator
			if token.Len() > 0 {
				result.WriteString(h.highlightGoToken(token.String(), keywords, types))
				token.Reset()
			}
			result.WriteString(h.operatorStyle.Render(string(r)))
			continue
		}

		// Handle whitespace
		if r == ' ' || r == '\t' {
			// If we have a token, highlight it before the whitespace
			if token.Len() > 0 {
				result.WriteString(h.highlightGoToken(token.String(), keywords, types))
				token.Reset()
			}
			result.WriteRune(r)
			continue
		}

		// Add to current token
		token.WriteRune(r)
	}

	// Handle any remaining token
	if token.Len() > 0 {
		result.WriteString(h.highlightGoToken(token.String(), keywords, types))
	}

	return result.String()
}

// highlightGoToken highlights a single Go token
func (h *SyntaxHighlighter) highlightGoToken(token string, keywords, types map[string]bool) string {
	// Check if token is a keyword
	if keywords[token] {
		return h.keywordStyle.Render(token)
	}

	// Check if token is a type
	if types[token] {
		return h.typeStyle.Render(token)
	}

	// Check if token is a number
	if isNumber(token) {
		return h.numberStyle.Render(token)
	}

	// Check if token is a function call
	if strings.HasSuffix(token, "(") {
		funcName := token[:len(token)-1]
		return h.functionStyle.Render(funcName) + h.operatorStyle.Render("(")
	}

	// Default style
	return h.defaultStyle.Render(token)
}

// HighlightPython highlights Python code
func (h *SyntaxHighlighter) HighlightPython(code string) string {
	lines := strings.Split(code, "\n")
	highlightedLines := make([]string, len(lines))

	// Keywords in Python
	keywords := map[string]bool{
		"def":      true,
		"class":    true,
		"import":   true,
		"from":     true,
		"as":       true,
		"if":       true,
		"elif":     true,
		"else":     true,
		"for":      true,
		"while":    true,
		"try":      true,
		"except":   true,
		"finally":  true,
		"with":     true,
		"return":   true,
		"break":    true,
		"continue": true,
		"pass":     true,
		"raise":    true,
		"assert":   true,
		"lambda":   true,
		"global":   true,
		"nonlocal": true,
		"True":     true,
		"False":    true,
		"None":     true,
		"and":      true,
		"or":       true,
		"not":      true,
		"is":       true,
		"in":       true,
		"yield":    true,
		"async":    true,
		"await":    true,
	}

	// Types in Python
	types := map[string]bool{
		"int":       true,
		"float":     true,
		"str":       true,
		"list":      true,
		"tuple":     true,
		"dict":      true,
		"set":       true,
		"bool":      true,
		"bytes":     true,
		"object":    true,
		"complex":   true,
		"range":     true,
		"Exception": true,
	}

	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			highlightedLines[i] = line
			continue
		}

		// Handle comments
		if commentIdx := strings.Index(line, "#"); commentIdx != -1 {
			beforeComment := line[:commentIdx]
			comment := line[commentIdx:]
			highlightedLines[i] = h.highlightPythonParts(beforeComment, keywords, types) + h.commentStyle.Render(comment)
			continue
		}

		highlightedLines[i] = h.highlightPythonParts(line, keywords, types)
	}

	return strings.Join(highlightedLines, "\n")
}

// highlightPythonParts highlights parts of a Python code line
func (h *SyntaxHighlighter) highlightPythonParts(line string, keywords, types map[string]bool) string {
	var result strings.Builder
	var token strings.Builder
	var inString bool
	var stringQuote rune

	for _, r := range line {
		// Handle strings
		if (r == '"' || r == '\'') && !inString {
			// If we have a token, highlight it before starting the string
			if token.Len() > 0 {
				result.WriteString(h.highlightPythonToken(token.String(), keywords, types))
				token.Reset()
			}
			inString = true
			stringQuote = r
			token.WriteRune(r)
			continue
		}

		if inString && r == stringQuote {
			// End of string
			token.WriteRune(r)
			result.WriteString(h.stringStyle.Render(token.String()))
			token.Reset()
			inString = false
			continue
		}

		if inString {
			token.WriteRune(r)
			continue
		}

		// Handle operators and delimiters
		if strings.ContainsRune("+-*/=<>!&|^~();:,.[]{}%", r) {
			// If we have a token, highlight it before the operator
			if token.Len() > 0 {
				result.WriteString(h.highlightPythonToken(token.String(), keywords, types))
				token.Reset()
			}
			result.WriteString(h.operatorStyle.Render(string(r)))
			continue
		}

		// Handle whitespace
		if r == ' ' || r == '\t' {
			// If we have a token, highlight it before the whitespace
			if token.Len() > 0 {
				result.WriteString(h.highlightPythonToken(token.String(), keywords, types))
				token.Reset()
			}
			result.WriteRune(r)
			continue
		}

		// Add to current token
		token.WriteRune(r)
	}

	// Handle any remaining token
	if token.Len() > 0 {
		result.WriteString(h.highlightPythonToken(token.String(), keywords, types))
	}

	return result.String()
}

// highlightPythonToken highlights a single Python token
func (h *SyntaxHighlighter) highlightPythonToken(token string, keywords, types map[string]bool) string {
	// Check if token is a keyword
	if keywords[token] {
		return h.keywordStyle.Render(token)
	}

	// Check if token is a type
	if types[token] {
		return h.typeStyle.Render(token)
	}

	// Check if token is a number
	if isNumber(token) {
		return h.numberStyle.Render(token)
	}

	// Check if token is a function call
	if strings.HasSuffix(token, "(") {
		funcName := token[:len(token)-1]
		return h.functionStyle.Render(funcName) + h.operatorStyle.Render("(")
	}

	// Default style
	return h.defaultStyle.Render(token)
}

// HighlightJavaScript highlights JavaScript code
func (h *SyntaxHighlighter) HighlightJavaScript(code string) string {
	lines := strings.Split(code, "\n")
	highlightedLines := make([]string, len(lines))

	// Keywords in JavaScript
	keywords := map[string]bool{
		"function":   true,
		"var":        true,
		"let":        true,
		"const":      true,
		"if":         true,
		"else":       true,
		"for":        true,
		"while":      true,
		"do":         true,
		"switch":     true,
		"case":       true,
		"default":    true,
		"break":      true,
		"continue":   true,
		"return":     true,
		"try":        true,
		"catch":      true,
		"finally":    true,
		"throw":      true,
		"class":      true,
		"extends":    true,
		"super":      true,
		"this":       true,
		"new":        true,
		"import":     true,
		"export":     true,
		"from":       true,
		"as":         true,
		"async":      true,
		"await":      true,
		"true":       true,
		"false":      true,
		"null":       true,
		"undefined":  true,
		"typeof":     true,
		"instanceof": true,
		"void":       true,
		"delete":     true,
	}

	// Types in JavaScript
	types := map[string]bool{
		"Array":    true,
		"Boolean":  true,
		"Date":     true,
		"Error":    true,
		"Function": true,
		"JSON":     true,
		"Math":     true,
		"Number":   true,
		"Object":   true,
		"Promise":  true,
		"RegExp":   true,
		"String":   true,
		"Symbol":   true,
		"Map":      true,
		"Set":      true,
		"WeakMap":  true,
		"WeakSet":  true,
	}

	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			highlightedLines[i] = line
			continue
		}

		// Handle comments
		if commentIdx := strings.Index(line, "//"); commentIdx != -1 {
			beforeComment := line[:commentIdx]
			comment := line[commentIdx:]
			highlightedLines[i] = h.highlightJavaScriptParts(beforeComment, keywords, types) + h.commentStyle.Render(comment)
			continue
		}

		highlightedLines[i] = h.highlightJavaScriptParts(line, keywords, types)
	}

	return strings.Join(highlightedLines, "\n")
}

// highlightJavaScriptParts highlights parts of a JavaScript code line
func (h *SyntaxHighlighter) highlightJavaScriptParts(line string, keywords, types map[string]bool) string {
	var result strings.Builder
	var token strings.Builder
	var inString bool
	var stringQuote rune

	for _, r := range line {
		// Handle strings
		if (r == '"' || r == '\'' || r == '`') && !inString {
			// If we have a token, highlight it before starting the string
			if token.Len() > 0 {
				result.WriteString(h.highlightJavaScriptToken(token.String(), keywords, types))
				token.Reset()
			}
			inString = true
			stringQuote = r
			token.WriteRune(r)
			continue
		}

		if inString && r == stringQuote {
			// End of string
			token.WriteRune(r)
			result.WriteString(h.stringStyle.Render(token.String()))
			token.Reset()
			inString = false
			continue
		}

		if inString {
			token.WriteRune(r)
			continue
		}

		// Handle operators and delimiters
		if strings.ContainsRune("+-*/=<>!&|^~();:,.[]{}%", r) {
			// If we have a token, highlight it before the operator
			if token.Len() > 0 {
				result.WriteString(h.highlightJavaScriptToken(token.String(), keywords, types))
				token.Reset()
			}
			result.WriteString(h.operatorStyle.Render(string(r)))
			continue
		}

		// Handle whitespace
		if r == ' ' || r == '\t' {
			// If we have a token, highlight it before the whitespace
			if token.Len() > 0 {
				result.WriteString(h.highlightJavaScriptToken(token.String(), keywords, types))
				token.Reset()
			}
			result.WriteRune(r)
			continue
		}

		// Add to current token
		token.WriteRune(r)
	}

	// Handle any remaining token
	if token.Len() > 0 {
		result.WriteString(h.highlightJavaScriptToken(token.String(), keywords, types))
	}

	return result.String()
}

// highlightJavaScriptToken highlights a single JavaScript token
func (h *SyntaxHighlighter) highlightJavaScriptToken(token string, keywords, types map[string]bool) string {
	// Check if token is a keyword
	if keywords[token] {
		return h.keywordStyle.Render(token)
	}

	// Check if token is a type
	if types[token] {
		return h.typeStyle.Render(token)
	}

	// Check if token is a number
	if isNumber(token) {
		return h.numberStyle.Render(token)
	}

	// Check if token is a function call
	if strings.HasSuffix(token, "(") {
		funcName := token[:len(token)-1]
		return h.functionStyle.Render(funcName) + h.operatorStyle.Render("(")
	}

	// Default style
	return h.defaultStyle.Render(token)
}

// HighlightCode highlights code based on the language
func (h *SyntaxHighlighter) HighlightCode(code, language string) string {
	switch strings.ToLower(language) {
	case "go":
		return h.HighlightGo(code)
	case "python", "py":
		return h.HighlightPython(code)
	case "javascript", "js":
		return h.HighlightJavaScript(code)
	default:
		return code
	}
}

// RenderCodeBlock renders a code block with syntax highlighting
func (h *SyntaxHighlighter) RenderCodeBlock(code, language string) string {
	// Highlight the code
	highlightedCode := h.HighlightCode(code, language)

	// Add a border and padding
	style := lipgloss.NewStyle().
		Background(lipgloss.Color(h.backgroundColor)).
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#404040"))

	// Add language label
	langLabel := fmt.Sprintf(" %s ", strings.ToUpper(language))
	langLabelStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#569CD6")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1)

	rendered := style.Render(highlightedCode)

	// Add language label at the top right
	lines := strings.Split(rendered, "\n")
	if len(lines) > 0 {
		labelLen := len(langLabel)
		lineLen := len(lines[0])
		if labelLen < lineLen {
			padding := strings.Repeat(" ", lineLen-labelLen-2)
			langLine := padding + langLabelStyle.Render(langLabel)
			lines[0] = langLine
		}
	}

	return strings.Join(lines, "\n")
}

// isNumber checks if a string represents a number
func isNumber(s string) bool {
	// Check for hex, octal, binary format
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0o") || strings.HasPrefix(s, "0b") {
		return true
	}

	// Check for decimal with possible decimal point
	hasDigit := false
	hasDecimal := false

	for i, r := range s {
		if r >= '0' && r <= '9' {
			hasDigit = true
		} else if r == '.' && !hasDecimal {
			hasDecimal = true
		} else if (r == 'e' || r == 'E') && hasDigit && i < len(s)-1 {
			// Scientific notation
			continue
		} else if (r == '+' || r == '-') && i > 0 && (s[i-1] == 'e' || s[i-1] == 'E') {
			// Sign in scientific notation
			continue
		} else {
			return false
		}
	}

	return hasDigit
}
