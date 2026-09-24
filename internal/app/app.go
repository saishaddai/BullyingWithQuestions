package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"bullyingwithquestions/internal/content"
	"bullyingwithquestions/internal/tui"
)

// Run starts QuickDeck using the default content directory.
func Run() error {
	return RunWithContentDir("content")
}

// LoadModel loads local content and builds the initial selection or error model.
func LoadModel(contentDir string) (tui.Model, error) {
	decks, warnings, err := content.LoadContentDir(contentDir)
	if err != nil {
		return tui.NewErrorModel(err.Error()), nil
	}
	model := tui.NewSelectionModel(decks, nil)
	model.SetWarnings(warnings)
	return model, nil
}

// RunWithContentDir loads local content and starts the selection screen.
func RunWithContentDir(contentDir string) error {
	model, err := LoadModel(contentDir)
	if err != nil {
		return err
	}
	program := tea.NewProgram(model.ProgramModel())
	_, err = program.Run()
	return err
}
