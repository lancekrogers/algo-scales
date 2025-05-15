// Package screens contains UI screens for different app states
package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lancekrogers/algo-scales/internal/problem"
	"github.com/lancekrogers/algo-scales/internal/ui/view"
)

// ProblemSelectionState tracks the current step in the problem selection process
type ProblemSelectionState int

const (
	StatePatternSelection ProblemSelectionState = iota
	StateProblemList
)

// ProblemSelectionModel represents the problem selection screen model
type ProblemSelectionModel struct {
	State             ProblemSelectionState
	Problems          []problem.Problem
	FilteredProblems  []problem.Problem
	Patterns          []string
	SelectedPattern   string
	SelectedProblemIdx int
	SelectedProblem   *problem.Problem
	Width             int
	Height            int
	Loading           bool
	Spinner           spinner.Model
	Language          string
	Mode              string
	Ready             bool
	PatternViz        *view.PatternVisualization
}

// NewProblemSelectionModel creates a new problem selection model
func NewProblemSelectionModel(allProblems []problem.Problem, language, mode string) ProblemSelectionModel {
	// Get all patterns
	patterns := problem.GetPatterns(allProblems)
	
	// Create spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
	
	return ProblemSelectionModel{
		State:             StatePatternSelection,
		Problems:          allProblems,
		FilteredProblems:  allProblems,
		Patterns:          patterns,
		SelectedPattern:   "",
		SelectedProblemIdx: 0,
		Loading:           false,
		Spinner:           s,
		Language:          language,
		Mode:              mode,
		PatternViz:        view.NewPatternVisualization(),
	}
}

// Init initializes the problem selection model
func (m ProblemSelectionModel) Init() tea.Cmd {
	return spinner.Tick
}

// Update handles updates to the problem selection model
func (m ProblemSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			switch m.State {
			case StatePatternSelection:
				// Set selected pattern and filter problems
				if m.SelectedProblemIdx == 0 {
					// "All Patterns" selected
					m.SelectedPattern = ""
					m.FilteredProblems = m.Problems
				} else if m.SelectedProblemIdx <= len(m.Patterns) {
					pattern := m.Patterns[m.SelectedProblemIdx-1]
					m.SelectedPattern = pattern
					m.FilteredProblems = problem.GetProblemsByPattern(m.Problems, pattern)
				}
				
				// Move to problem list state
				m.State = StateProblemList
				m.SelectedProblemIdx = 0
				
			case StateProblemList:
				// Set selected problem if valid
				if m.SelectedProblemIdx >= 0 && m.SelectedProblemIdx < len(m.FilteredProblems) {
					m.SelectedProblem = &m.FilteredProblems[m.SelectedProblemIdx]
					// Return selected problem to caller
					return m, func() tea.Msg {
						return problemSelectedMsg{
							Problem: m.SelectedProblem,
							Pattern: m.SelectedPattern,
						}
					}
				}
			}

		case "esc", "backspace":
			if m.State == StateProblemList {
				// Go back to pattern selection
				m.State = StatePatternSelection
				m.SelectedProblemIdx = 0
			} else {
				// Exit if already in pattern selection
				return m, tea.Quit
			}

		case "up", "k":
			// Move selection up
			if m.SelectedProblemIdx > 0 {
				m.SelectedProblemIdx--
			}

		case "down", "j":
			// Move selection down
			switch m.State {
			case StatePatternSelection:
				// +1 for "All Patterns" option
				if m.SelectedProblemIdx < len(m.Patterns) {
					m.SelectedProblemIdx++
				}
			case StateProblemList:
				if m.SelectedProblemIdx < len(m.FilteredProblems)-1 {
					m.SelectedProblemIdx++
				}
			}
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Ready = true

	case spinner.TickMsg:
		var spinnerCmd tea.Cmd
		m.Spinner, spinnerCmd = m.Spinner.Update(msg)
		return m, spinnerCmd
	}

	return m, cmd
}

// View renders the problem selection screen
func (m ProblemSelectionModel) View() string {
	if !m.Ready {
		return "Loading..."
	}

	var content string

	switch m.State {
	case StatePatternSelection:
		content = m.renderPatternSelection()
	case StateProblemList:
		content = m.renderProblemList()
	}

	// Add loading indicator if loading
	if m.Loading {
		content += "\n\n" + m.Spinner.View() + " Loading..."
	}

	// Add navigation help
	navigationHelp := "↑/↓: Navigate • Enter: Select • Backspace: Back • q: Quit"
	content += "\n\n" + view.HelpStyle.Render(navigationHelp)

	// Center the content
	return content
}

// renderPatternSelection renders the pattern selection screen
func (m ProblemSelectionModel) renderPatternSelection() string {
	title := view.TitleStyle.Render("Select Algorithm Pattern")
	subtitle := view.SubtitleStyle.Render("Choose a pattern to focus on or select 'All Patterns'")

	// Create pattern options
	var options strings.Builder
	
	// Add "All Patterns" option
	allPatternsOption := ""
	if m.SelectedProblemIdx == 0 {
		allPatternsOption = view.FocusedItemStyle.Render("▶ All Patterns")
	} else {
		allPatternsOption = view.UnfocusedItemStyle.Render("  All Patterns")
	}
	options.WriteString(allPatternsOption + "\n\n")

	// Add each pattern with description and color based on musical scale
	for i, pattern := range m.Patterns {
		option := ""
		
		// Get pattern information
		scale, ok := view.MusicScales[pattern]
		if !ok {
			// Skip patterns without visualization information
			continue
		}
		
		// Create style based on pattern
		patternStyle, _, _ := view.GetPatternStyle(pattern)
		
		// Format the option
		if i+1 == m.SelectedProblemIdx {
			option = view.FocusedItemStyle.Render(fmt.Sprintf("▶ %s", scale.Name))
		} else {
			option = patternStyle.Render(fmt.Sprintf("  %s", scale.Name))
		}
		
		// Add description
		options.WriteString(fmt.Sprintf("%s\n   %s\n\n", option, scale.Description))
	}

	// Add pattern visualization for selected pattern
	var visualization string
	if m.SelectedProblemIdx > 0 && m.SelectedProblemIdx <= len(m.Patterns) {
		selectedPattern := m.Patterns[m.SelectedProblemIdx-1]
		visualization = "\n" + view.BorderedBoxStyle.Render(
			view.HeaderStyle.Render("Pattern Visualization") + "\n\n" +
			m.PatternViz.VisualizePattern(selectedPattern, "", m.Width-20),
		)
	}

	return title + "\n\n" + 
		subtitle + "\n\n" + 
		view.MenuBoxStyle.Render(options.String()) + 
		visualization
}

// renderProblemList renders the problem list screen
func (m ProblemSelectionModel) renderProblemList() string {
	// Create title with pattern information
	var title string
	if m.SelectedPattern != "" {
		if scale, ok := view.MusicScales[m.SelectedPattern]; ok {
			title = view.TitleStyle.Render(fmt.Sprintf("Select Problem - %s", scale.Name))
		} else {
			title = view.TitleStyle.Render("Select Problem")
		}
	} else {
		title = view.TitleStyle.Render("Select Problem")
	}
	
	// Create subtitle with count
	subtitle := view.SubtitleStyle.Render(
		fmt.Sprintf("Found %d problems - Select one to start practicing", len(m.FilteredProblems)),
	)

	// Handle no problems case
	if len(m.FilteredProblems) == 0 {
		return title + "\n\n" + 
			subtitle + "\n\n" + 
			"No problems found for this pattern."
	}

	// Create problem list
	var problemList strings.Builder
	for i, prob := range m.FilteredProblems {
		// Create difficulty style
		var difficultyStyle lipgloss.Style
		switch prob.Difficulty {
		case "easy":
			difficultyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#2ecc71"))
		case "medium":
			difficultyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f1c40f"))
		case "hard":
			difficultyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#e74c3c"))
		default:
			difficultyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7f8c8d"))
		}
		
		// Format option
		option := ""
		if i == m.SelectedProblemIdx {
			option = view.FocusedItemStyle.Render(fmt.Sprintf("▶ %s", prob.Title))
		} else {
			option = view.UnfocusedItemStyle.Render(fmt.Sprintf("  %s", prob.Title))
		}
		
		// Format difficulty and time
		difficulty := difficultyStyle.Render(strings.Title(prob.Difficulty))
		timeEstimate := fmt.Sprintf("%d min", prob.EstimatedTime)
		
		problemList.WriteString(fmt.Sprintf("%-40s [%s | %s]\n\n", option, difficulty, timeEstimate))
	}

	// Add problem preview for selected problem
	var preview string
	if m.SelectedProblemIdx >= 0 && m.SelectedProblemIdx < len(m.FilteredProblems) {
		selectedProblem := m.FilteredProblems[m.SelectedProblemIdx]
		
		preview = "\n" + view.BorderedBoxStyle.Render(
			view.HeaderStyle.Render("Problem Preview") + "\n\n" +
			selectedProblem.Description + "\n\n" +
			view.HeaderStyle.Render("First Example") + "\n\n" +
			fmt.Sprintf("Input: %s\n", selectedProblem.Examples[0].Input) +
			fmt.Sprintf("Output: %s", selectedProblem.Examples[0].Output),
		)
	}

	return title + "\n\n" + 
		subtitle + "\n\n" + 
		view.MenuBoxStyle.Render(problemList.String()) + 
		preview
}

// Custom message for when a problem is selected
type problemSelectedMsg struct {
	Problem *problem.Problem
	Pattern string
}