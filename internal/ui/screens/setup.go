// Package screens contains UI screens for different app states
package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lancekrogers/algo-scales/internal/common/config"
	"github.com/lancekrogers/algo-scales/internal/ui/view"
)

// SetupScreenState tracks the current step in the setup process
type SetupScreenState int

const (
	StateLanguage SetupScreenState = iota
	StateTimer
	StateMode
	StateConfirmation
)

// SetupModel represents the setup screen model
type SetupModel struct {
	State           SetupScreenState
	UserConfig      config.UserConfig
	LanguageOptions []string
	TimerOptions    []int
	ModeOptions     []string
	SelectedIndex   int
	textInput       textinput.Model
	width           int
	height          int
	successMsg      string
	errorMsg        string
}

// NewSetupModel creates a new setup screen model
func NewSetupModel() SetupModel {
	// Load user config if exists
	cfg, err := config.LoadConfig()
	if err != nil {
		cfg = config.DefaultConfig()
	}

	// Set up options
	languageOptions := config.ListLanguages()
	timerOptions := config.ListTimerOptions()
	modeOptions := config.ListModes()

	// Set up text input
	ti := textinput.New()
	ti.Placeholder = "Type custom time (minutes)..."
	ti.CharLimit = 3
	ti.Width = 20

	return SetupModel{
		State:           StateLanguage,
		UserConfig:      cfg,
		LanguageOptions: languageOptions,
		TimerOptions:    timerOptions,
		ModeOptions:     modeOptions,
		SelectedIndex:   0,
		textInput:       ti,
	}
}

// Init initializes the setup screen
func (m SetupModel) Init() tea.Cmd {
	return nil
}

// Update handles updates to the setup screen
func (m SetupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			// Exit without saving
			return m, tea.Quit

		case "enter":
			return m.handleEnter()

		case "up", "k":
			// Move selection up
			if m.SelectedIndex > 0 {
				m.SelectedIndex--
			}

		case "down", "j":
			// Move selection down
			switch m.State {
			case StateLanguage:
				if m.SelectedIndex < len(m.LanguageOptions)-1 {
					m.SelectedIndex++
				}
			case StateTimer:
				if m.SelectedIndex < len(m.TimerOptions) {
					m.SelectedIndex++
				}
			case StateMode:
				if m.SelectedIndex < len(m.ModeOptions)-1 {
					m.SelectedIndex++
				}
			}

		case "tab":
			// Move to next screen
			return m.moveToNextState()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	// Handle text input if in custom timer mode
	if m.State == StateTimer && m.SelectedIndex == len(m.TimerOptions) {
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

// handleEnter handles the Enter key press
func (m SetupModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.State {
	case StateLanguage:
		// Set selected language
		m.UserConfig.Language = m.LanguageOptions[m.SelectedIndex]
		m.State = StateTimer
		m.SelectedIndex = 0
	case StateTimer:
		// Set selected timer
		if m.SelectedIndex < len(m.TimerOptions) {
			m.UserConfig.TimerDuration = m.TimerOptions[m.SelectedIndex]
		} else {
			// Parse custom timer input
			var customTime int
			_, err := fmt.Sscanf(m.textInput.Value(), "%d", &customTime)
			if err == nil && customTime > 0 {
				m.UserConfig.TimerDuration = customTime
			} else {
				// Invalid input, use default
				m.UserConfig.TimerDuration = config.DefaultConfig().TimerDuration
			}
		}
		m.State = StateMode
		m.SelectedIndex = 0
	case StateMode:
		// Set selected mode
		m.UserConfig.Mode = m.ModeOptions[m.SelectedIndex]
		m.State = StateConfirmation
		m.SelectedIndex = 0
	case StateConfirmation:
		// Save config or start practice based on selection
		if m.SelectedIndex == 0 {
			// Save config
			err := config.SaveConfig(m.UserConfig)
			if err != nil {
				m.errorMsg = fmt.Sprintf("Failed to save config: %v", err)
			} else {
				m.successMsg = "Configuration saved successfully!"
			}
		}
		// Start practice (will be handled by the caller)
		return m, tea.Quit
	}

	return m, nil
}

// moveToNextState moves to the next setup state
func (m SetupModel) moveToNextState() (tea.Model, tea.Cmd) {
	switch m.State {
	case StateLanguage:
		m.UserConfig.Language = m.LanguageOptions[m.SelectedIndex]
		m.State = StateTimer
		m.SelectedIndex = 0
	case StateTimer:
		if m.SelectedIndex < len(m.TimerOptions) {
			m.UserConfig.TimerDuration = m.TimerOptions[m.SelectedIndex]
		} else {
			// Parse custom timer input
			var customTime int
			_, err := fmt.Sscanf(m.textInput.Value(), "%d", &customTime)
			if err == nil && customTime > 0 {
				m.UserConfig.TimerDuration = customTime
			} else {
				// Invalid input, use default
				m.UserConfig.TimerDuration = config.DefaultConfig().TimerDuration
			}
		}
		m.State = StateMode
		m.SelectedIndex = 0
	case StateMode:
		m.UserConfig.Mode = m.ModeOptions[m.SelectedIndex]
		m.State = StateConfirmation
		m.SelectedIndex = 0
	case StateConfirmation:
		// Do nothing, wait for Enter
	}

	return m, nil
}

// View renders the setup screen
func (m SetupModel) View() string {
	var content string

	switch m.State {
	case StateLanguage:
		content = m.renderLanguageSelection()
	case StateTimer:
		content = m.renderTimerSelection()
	case StateMode:
		content = m.renderModeSelection()
	case StateConfirmation:
		content = m.renderConfirmation()
	}

	// Add error or success message
	if m.errorMsg != "" {
		content += "\n\n" + view.ErrorStyle.Render(m.errorMsg)
	}
	if m.successMsg != "" {
		content += "\n\n" + view.SuccessStyle.Render(m.successMsg)
	}

	// Add navigation help
	navigationHelp := "↑/↓: Navigate • Enter: Select • Tab: Next • Esc: Quit"
	content += "\n\n" + view.HelpStyle.Render(navigationHelp)

	// Center the content
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

// renderLanguageSelection renders the language selection screen
func (m SetupModel) renderLanguageSelection() string {
	title := view.TitleStyle.Render("Select Programming Language")
	subtitle := view.SubtitleStyle.Render("Choose your preferred language for code examples and solutions")

	// Render language options
	var options strings.Builder
	for i, lang := range m.LanguageOptions {
		option := ""
		if i == m.SelectedIndex {
			option = view.FocusedItemStyle.Render(fmt.Sprintf("▶ %s", strings.Title(lang)))
		} else {
			option = view.UnfocusedItemStyle.Render(fmt.Sprintf("  %s", strings.Title(lang)))
		}
		options.WriteString(option + "\n")
	}

	return view.MenuBoxStyle.Render(
		title + "\n\n" + 
		subtitle + "\n\n" + 
		options.String(),
	)
}

// renderTimerSelection renders the timer selection screen
func (m SetupModel) renderTimerSelection() string {
	title := view.TitleStyle.Render("Configure Timer")
	subtitle := view.SubtitleStyle.Render("Set your preferred time limit for solving problems")

	// Render timer options
	var options strings.Builder
	for i, timer := range m.TimerOptions {
		option := ""
		if i == m.SelectedIndex {
			option = view.FocusedItemStyle.Render(fmt.Sprintf("▶ %d minutes", timer))
		} else {
			option = view.UnfocusedItemStyle.Render(fmt.Sprintf("  %d minutes", timer))
		}
		options.WriteString(option + "\n")
	}

	// Add custom timer option
	customOption := ""
	if m.SelectedIndex == len(m.TimerOptions) {
		customOption = view.FocusedItemStyle.Render("▶ Custom: ") + m.textInput.View()
	} else {
		customOption = view.UnfocusedItemStyle.Render("  Custom: ") + m.textInput.View()
	}
	options.WriteString(customOption)

	return view.MenuBoxStyle.Render(
		title + "\n\n" + 
		subtitle + "\n\n" + 
		options.String(),
	)
}

// renderModeSelection renders the mode selection screen
func (m SetupModel) renderModeSelection() string {
	title := view.TitleStyle.Render("Select Learning Mode")
	subtitle := view.SubtitleStyle.Render("Choose how you want to learn algorithm patterns")

	// Render mode options with descriptions
	var options strings.Builder
	modeDescriptions := []string{
		"Detailed explanations and visible solutions",
		"Practice with hints available",
		"Rapid-fire practice with strict time limits",
	}

	for i, mode := range m.ModeOptions {
		option := ""
		if i == m.SelectedIndex {
			option = view.FocusedItemStyle.Render(fmt.Sprintf("▶ %s Mode", strings.Title(mode)))
		} else {
			option = view.UnfocusedItemStyle.Render(fmt.Sprintf("  %s Mode", strings.Title(mode)))
		}
		options.WriteString(fmt.Sprintf("%s\n   %s\n\n", option, modeDescriptions[i]))
	}

	return view.MenuBoxStyle.Render(
		title + "\n\n" + 
		subtitle + "\n\n" + 
		options.String(),
	)
}

// renderConfirmation renders the confirmation screen
func (m SetupModel) renderConfirmation() string {
	title := view.TitleStyle.Render("Configuration Complete")

	// Show selected settings
	settings := fmt.Sprintf(
		"Language: %s\nTimer: %d minutes\nMode: %s",
		strings.Title(m.UserConfig.Language),
		m.UserConfig.TimerDuration,
		strings.Title(m.UserConfig.Mode),
	)
	settingsBox := view.BorderedBoxStyle.Render(settings)

	// Render confirmation options
	var options strings.Builder
	confirmOptions := []string{
		"Save configuration and continue",
		"Start practice without saving",
	}

	for i, option := range confirmOptions {
		if i == m.SelectedIndex {
			options.WriteString(view.FocusedItemStyle.Render(fmt.Sprintf("▶ %s", option)) + "\n")
		} else {
			options.WriteString(view.UnfocusedItemStyle.Render(fmt.Sprintf("  %s", option)) + "\n")
		}
	}

	return view.MenuBoxStyle.Render(
		title + "\n\n" + 
		settingsBox + "\n\n" +
		"What would you like to do?\n\n" +
		options.String(),
	)
}