package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// KeyMap defines all keyboard shortcuts for the application
type KeyMap struct {
	// Navigation
	Up        key.Binding
	Down      key.Binding
	Left      key.Binding
	Right     key.Binding
	PageUp    key.Binding
	PageDown  key.Binding
	Home      key.Binding
	End       key.Binding
	
	// Actions
	Select    key.Binding
	Back      key.Binding
	Quit      key.Binding
	Help      key.Binding
	Refresh   key.Binding
	
	// Session specific
	Edit      key.Binding
	Test      key.Binding
	Hint      key.Binding
	Solution  key.Binding
	Pause     key.Binding
	Submit    key.Binding
	
	// List specific
	Filter    key.Binding
	Sort      key.Binding
	Search    key.Binding
	
	// Settings specific
	Save      key.Binding
	Cancel    key.Binding
	Reset     key.Binding
	
	// Daily specific
	Next      key.Binding
	Previous  key.Binding
	Skip      key.Binding
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		// Navigation
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "ctrl+u"),
			key.WithHelp("pgup/ctrl+u", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", "ctrl+d"),
			key.WithHelp("pgdn/ctrl+d", "page down"),
		),
		Home: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("home/g", "go to start"),
		),
		End: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("end/G", "go to end"),
		),
		
		// Actions
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc", "backspace"),
			key.WithHelp("esc", "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r", "ctrl+r"),
			key.WithHelp("r", "refresh"),
		),
		
		// Session specific
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit code"),
		),
		Test: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "run tests"),
		),
		Hint: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "toggle hint"),
		),
		Solution: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "show solution"),
		),
		Pause: key.NewBinding(
			key.WithKeys("p", "space"),
			key.WithHelp("p/space", "pause timer"),
		),
		Submit: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "submit solution"),
		),
		
		// List specific
		Filter: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "filter"),
		),
		Sort: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "sort"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		
		// Settings specific
		Save: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "save"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		Reset: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("ctrl+r", "reset"),
		),
		
		// Daily specific
		Next: key.NewBinding(
			key.WithKeys("n", "tab"),
			key.WithHelp("n/tab", "next"),
		),
		Previous: key.NewBinding(
			key.WithKeys("N", "shift+tab"),
			key.WithHelp("N/shift+tab", "previous"),
		),
		Skip: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "skip"),
		),
	}
}

// FullHelp returns all help items
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		// Navigation
		{k.Up, k.Down, k.Left, k.Right},
		// Actions
		{k.Select, k.Back, k.Quit, k.Help},
		// Page navigation
		{k.PageUp, k.PageDown, k.Home, k.End},
	}
}

// ShortHelp returns essential help items
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// Update keymap in model
func (m Model) updateKeys() Model {
	// Global keys are always active
	m.keys = globalKeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q/ctrl+c", "quit"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc", "left"),
			key.WithHelp("esc/←", "back"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
	}
	
	return m
}

// Helper to check key matches
func matchesKey(msg tea.KeyMsg, binding key.Binding) bool {
	return key.Matches(msg, binding)
}