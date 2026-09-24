package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"bullyingwithquestions/internal/content"
)

type screen int

const (
	selectionScreen screen = iota
	errorScreen
)

// DeckSelectedMsg is emitted when the user confirms a deck.
type DeckSelectedMsg struct {
	Deck content.LoadedDeck
}

// Model is the Phase 4 deck-selection model.
type Model struct {
	decks    []content.LoadedDeck
	selected int
	offset   int
	width    int
	height   int
	screen   screen
	err      string
	warnings []string
}

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F4A261"))
	mutedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#8A9BA8"))
	selectedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#2A9D8F"))
	panelStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2)
)

// NewSelectionModel creates a deck-selection model. The random argument is
// reserved for later session creation and keeps the constructor stable.
func NewSelectionModel(decks []content.LoadedDeck, _ interface{}) Model {
	return Model{decks: append([]content.LoadedDeck(nil), decks...), screen: selectionScreen}
}

// NewErrorModel creates a model that displays a startup or loading error.
func NewErrorModel(message string) Model {
	return Model{screen: errorScreen, err: message}
}

// SetWarnings adds non-fatal content warnings to the selection screen.
func (model *Model) SetWarnings(warnings []string) {
	model.warnings = append([]string(nil), warnings...)
}

// Init implements the Bubble Tea initialization contract.
func (model Model) Init() tea.Cmd {
	return nil
}

// Update handles keyboard navigation, selection, quit, and terminal resizing.
func (model Model) Update(message tea.Msg) (Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		model.width, model.height = message.Width, message.Height
		model.ensureVisible()
		return model, nil
	case tea.KeyMsg:
		if message.Type == tea.KeyCtrlC || message.Type == tea.KeyEscape || message.String() == "q" {
			return model, tea.Quit
		}
		if model.screen != selectionScreen || len(model.decks) == 0 {
			return model, nil
		}
		switch message.Type {
		case tea.KeyUp:
			if model.selected > 0 {
				model.selected--
				model.ensureVisible()
			}
		case tea.KeyDown:
			if model.selected < len(model.decks)-1 {
				model.selected++
				model.ensureVisible()
			}
		case tea.KeyEnter:
			return model, func() tea.Msg { return DeckSelectedMsg{Deck: model.decks[model.selected]} }
		}
	}
	return model, nil
}

// View renders the selection or error screen.
func (model Model) View() string {
	if model.screen == errorScreen {
		return panelStyle.Render(titleStyle.Render("QuickDeck") + "\n\n" + model.err + "\n\n" + mutedStyle.Render("q Quit"))
	}
	if model.width > 0 && model.width < 40 {
		return titleStyle.Render("QuickDeck") + "\n\n" + "Terminal too small. Resize to at least 40 columns."
	}

	header := titleStyle.Render("QuickDeck") + "  " + mutedStyle.Render("Choose a deck")
	if len(model.decks) == 0 {
		return panelStyle.Render(header + "\n\nNo decks available.\n\n" + mutedStyle.Render("q Quit"))
	}

	rows := make([]string, 0, model.visibleCount())
	end := model.offset + model.visibleCount()
	if end > len(model.decks) {
		end = len(model.decks)
	}
	for index := model.offset; index < end; index++ {
		prefix := "[ ] "
		row := model.decks[index].Name
		if index == model.selected {
			prefix = "[>] "
			row = selectedStyle.Render(row)
		}
		rows = append(rows, prefix+row)
	}

	description := model.decks[model.selected].Description
	body := strings.Join(rows, "\n") + "\n\n" + mutedStyle.Render(description)
	if len(model.warnings) > 0 {
		body += "\n\n" + mutedStyle.Render(fmt.Sprintf("%d content warning(s)", len(model.warnings)))
	}
	footer := mutedStyle.Render("Up/Down Move  Enter Select  q Quit")
	return panelStyle.Render(header + "\n\n" + body + "\n\n" + footer)
}

// SelectedIndex returns the highlighted deck index.
func (model Model) SelectedIndex() int {
	return model.selected
}

// ProgramModel adapts the testable concrete model to tea.Model.
func (model Model) ProgramModel() tea.Model {
	return programModel{model: model}
}

type programModel struct{ model Model }

func (model programModel) Init() tea.Cmd { return model.model.Init() }

func (model programModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	updated, command := model.model.Update(message)
	return programModel{model: updated}, command
}

func (model programModel) View() string { return model.model.View() }

func (model Model) visibleCount() int {
	count := model.height - 10
	if count < 1 {
		count = len(model.decks)
	}
	return count
}

func (model *Model) ensureVisible() {
	visible := model.visibleCount()
	if model.selected < model.offset {
		model.offset = model.selected
	}
	if model.selected >= model.offset+visible {
		model.offset = model.selected - visible + 1
	}
	if model.offset < 0 {
		model.offset = 0
	}
}

func (model Model) String() string {
	return fmt.Sprintf("selection screen, %d decks", len(model.decks))
}
