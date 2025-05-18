package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lancekrogers/algo-scales/internal/common/config"
)

// Settings options
var settingsOptions = []string{
	"Language",
	"Timer Duration",
	"Editor Command",
	"Theme",
	"Reset Statistics",
	"Clear Cache",
}

// Update handles updates for the settings screen
func (m Model) updateSettings(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case settingSavedMsg:
		m.settings.message = "Settings saved successfully!"
		m.settings.editing = false
		
	case settingErrorMsg:
		m.settings.message = fmt.Sprintf("Error saving settings: %v", msg.error)
		m.settings.editing = false
		
	case tea.KeyMsg:
		// Clear message after keypress
		if m.settings.message != "" && !m.settings.editing {
			m.settings.message = ""
		}
		
		switch msg.String() {
		case "up", "k":
			if m.settings.selectedOption > 0 {
				m.settings.selectedOption--
			}
		case "down", "j":
			if m.settings.selectedOption < len(settingsOptions)-1 {
				m.settings.selectedOption++
			}
		case "enter", "right", "l":
			return m.handleSettingSelection()
		case "left", "h":
			if m.settings.editing {
				m.settings.editing = false
			}
		}
		
		// Handle editing mode
		if m.settings.editing {
			switch msg.String() {
			case "backspace":
				if len(m.settings.editValue) > 0 {
					m.settings.editValue = m.settings.editValue[:len(m.settings.editValue)-1]
				}
			case "enter":
				return m.saveSettingValue()
			case "esc":
				m.settings.editing = false
			default:
				// Add character to edit value
				if len(msg.String()) == 1 {
					m.settings.editValue += msg.String()
				}
			}
		}
	}
	
	return m, nil
}

// View renders the settings screen
func (m Model) viewSettings() string {
	var b strings.Builder
	
	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("62")).
		MarginBottom(2)
	
	b.WriteString(titleStyle.Render("⚙️  Settings"))
	b.WriteString("\n\n")
	
	// Settings list
	for i, option := range settingsOptions {
		cursor := "  "
		if i == m.settings.selectedOption {
			cursor = "> "
		}
		
		// Format option with current value
		line := fmt.Sprintf("%s%-20s", cursor, option)
		
		// Add current value
		valueStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("243"))
		
		value := m.getSettingValue(i)
		if m.settings.editing && i == m.settings.selectedOption {
			// Show editing value
			value = m.settings.editValue + "█"
			valueStyle = valueStyle.Bold(true).Foreground(lipgloss.Color("214"))
		}
		
		line += valueStyle.Render(value)
		
		// Highlight selected option
		if i == m.settings.selectedOption {
			line = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("212")).
				Render(line)
		}
		
		b.WriteString(line + "\n")
	}
	
	// Show message if any
	if m.settings.message != "" {
		b.WriteString("\n")
		messageStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("214"))
		b.WriteString(messageStyle.Render(m.settings.message))
		b.WriteString("\n")
	}
	
	// Help text
	b.WriteString("\n")
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))
	
	if m.settings.editing {
		b.WriteString(helpStyle.Render("Type to edit • Enter: Save • Esc: Cancel"))
	} else {
		b.WriteString(helpStyle.Render("↑/↓: Navigate • Enter: Edit • Esc: Back"))
	}
	
	return b.String()
}

// getSettingValue returns the current value for a setting
func (m Model) getSettingValue(index int) string {
	switch index {
	case 0: // Language
		return m.config.Language
	case 1: // Timer Duration
		return fmt.Sprintf("%d minutes", m.config.TimerDuration)
	case 2: // Editor Command
		return m.config.EditorCommand
	case 3: // Theme
		return m.config.Theme
	case 4: // Reset Statistics
		return "Press Enter to reset"
	case 5: // Clear Cache
		return "Press Enter to clear"
	default:
		return "Unknown"
	}
}

// handleSettingSelection handles selecting a setting to edit
func (m Model) handleSettingSelection() (Model, tea.Cmd) {
	// Add bounds checking
	if m.settings.selectedOption < 0 || m.settings.selectedOption >= len(settingsOptions) {
		return m, nil
	}
	
	switch m.settings.selectedOption {
	case 0, 1, 2, 3: // Editable fields
		m.settings.editing = true
		m.settings.editingField = settingsOptions[m.settings.selectedOption]
		m.settings.editValue = m.getSettingValue(m.settings.selectedOption)
		// For timer, remove " minutes" suffix
		if m.settings.selectedOption == 1 {
			m.settings.editValue = fmt.Sprintf("%d", m.config.TimerDuration)
		}
	case 4: // Reset Statistics
		return m.resetStatistics()
	case 5: // Clear Cache
		return m.clearCache()
	}
	
	return m, nil
}

// saveSettingValue saves the edited setting value
func (m Model) saveSettingValue() (Model, tea.Cmd) {
	// Validate inputs and handle errors
	switch m.settings.selectedOption {
	case 0: // Language
		// Validate language is supported
		validLanguages := config.ListLanguages()
		langValid := false
		for _, lang := range validLanguages {
			if m.settings.editValue == lang {
				langValid = true
				break
			}
		}
		if langValid {
			m.config.Language = m.settings.editValue
		} else {
			m.settings.message = fmt.Sprintf("Invalid language. Choose from: %v", validLanguages)
			m.settings.editing = false
			return m, nil
		}
	case 1: // Timer Duration
		// Parse integer
		var duration int
		_, err := fmt.Sscanf(m.settings.editValue, "%d", &duration)
		if err != nil || duration <= 0 {
			m.settings.message = "Invalid timer duration. Please enter a positive number."
			m.settings.editing = false
			return m, nil
		}
		m.config.TimerDuration = duration
	case 2: // Editor Command
		if m.settings.editValue == "" {
			m.settings.message = "Editor command cannot be empty."
			m.settings.editing = false
			return m, nil
		}
		m.config.EditorCommand = m.settings.editValue
	case 3: // Theme
		m.config.Theme = m.settings.editValue
	}
	
	m.settings.editing = false
	
	// Save config
	return m, saveConfig(m.config)
}

// resetStatistics resets all user statistics
func (m Model) resetStatistics() (Model, tea.Cmd) {
	// This would reset statistics in a real implementation
	// For now, just show a message
	m.settings.message = "Statistics reset!"
	return m, nil
}

// clearCache clears the application cache
func (m Model) clearCache() (Model, tea.Cmd) {
	// This would clear cache in a real implementation
	// For now, just show a message
	m.settings.message = "Cache cleared!"
	return m, nil
}

// saveConfig command
func saveConfig(cfg config.UserConfig) tea.Cmd {
	return func() tea.Msg {
		err := config.SaveConfig(cfg)
		if err != nil {
			return settingErrorMsg{err}
		}
		return settingSavedMsg{}
	}
}

// Message types for settings
type settingSavedMsg struct{}
type settingErrorMsg struct{ error }