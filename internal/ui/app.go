// Package ui implements the terminal user interface for Algo Scales
package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lancekrogers/algo-scales/internal/ui/controller"
	"github.com/lancekrogers/algo-scales/internal/ui/model"
	"github.com/lancekrogers/algo-scales/internal/ui/view"
)

// App represents the UI application
type App struct {
	model      *model.UIModel
	view       *view.View
	controller *controller.Controller
}

// NewApp creates a new UI application
func NewApp() *App {
	// Create model
	m := model.NewModel()

	// Create view and controller
	v := view.NewView(&m)
	c := controller.NewController(&m)

	return &App{
		model:      &m,
		view:       v,
		controller: c,
	}
}

// Start launches the UI application
func (a *App) Start() error {
	// Create the Bubble Tea program
	p := tea.NewProgram(a)
	
	// Run the program
	_, err := p.Run()
	return err
}

// Implement tea.Model interface for the App

// Init initializes the application
func (a *App) Init() tea.Cmd {
	return a.controller.Init()
}

// Update handles updates to the application
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return a.controller.Update(msg)
}

// View renders the application
func (a *App) View() string {
	return a.view.Render()
}

// Messages that can be passed between components

// Define message types for bubbletea
type (
	// errorMsg is used to indicate an error
	errorMsg string

	// problemsLoadedMsg is sent when the problem list is loaded
	problemsLoadedMsg struct {
		Problems []model.Problem
	}

	// tickMsg is sent every second for timers
	tickMsg struct{}

	// problemSelectedMsg is sent when a problem is selected
	problemSelectedMsg struct {
		ProblemID string
		Mode      string
	}

	// codeUpdatedMsg is sent when the code is updated in the editor
	codeUpdatedMsg struct{}

	// testResultsMsg is sent after running tests
	testResultsMsg struct {
		Results   []model.TestResult
		AllPassed bool
	}

	// achievementUnlockedMsg is sent when an achievement is unlocked
	achievementUnlockedMsg struct {
		AchievementID string
	}
)