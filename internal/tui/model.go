package tui

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"bullyingwithquestions/internal/content"
	"bullyingwithquestions/internal/session"
)

type screen int

const (
	selectionScreen screen = iota
	studyScreen
	summaryScreen
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
	deck     content.LoadedDeck
	study    *session.Session
	random   *rand.Rand
	err      string
	warnings []string
}

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F4A261"))
	mutedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#8A9BA8"))
	selectedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#2A9D8F"))
	panelStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2)
)

// NewSelectionModel creates a deck-selection model with an optional session-local random source.
func NewSelectionModel(decks []content.LoadedDeck, source *rand.Rand) Model {
	return Model{decks: append([]content.LoadedDeck(nil), decks...), screen: selectionScreen, random: source}
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
	case DeckSelectedMsg:
		study, err := session.New(message.Deck.Cards, model.random, time.Now())
		if err != nil {
			model.screen = errorScreen
			model.err = fmt.Sprintf("deck %q: %v", message.Deck.ID, err)
			return model, nil
		}
		model.deck = message.Deck
		model.study = study
		model.screen = studyScreen
		return model, nil
	case tea.WindowSizeMsg:
		model.width, model.height = message.Width, message.Height
		model.ensureVisible()
		return model, nil
	case tea.KeyMsg:
		if message.Type == tea.KeyCtrlC || message.Type == tea.KeyEscape || message.String() == "q" {
			return model, tea.Quit
		}
		switch model.screen {
		case selectionScreen:
			if len(model.decks) == 0 {
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
		case studyScreen:
			if model.study == nil {
				return model, nil
			}
			switch message.Type {
			case tea.KeyLeft:
				model.study.Previous()
			case tea.KeyRight:
				if !model.study.Next() {
					model.screen = summaryScreen
				}
			case tea.KeySpace:
				model.study.Reveal()
			default:
				switch message.String() {
				case "h":
					model.study.Previous()
				case "l":
					if !model.study.Next() {
						model.screen = summaryScreen
					}
				case " ":
					model.study.Reveal()
				}
			}
		case summaryScreen:
			if message.Type == tea.KeyEnter {
				model.screen = selectionScreen
				model.study = nil
			}
		}
	}
	return model, nil
}

// View renders the current application screen.
func (model Model) View() string {
	if model.screen == errorScreen {
		return panelStyle.Render(titleStyle.Render("QuickDeck") + "\n\n" + model.err + "\n\n" + mutedStyle.Render("q Quit"))
	}
	if model.width > 0 && model.width < 40 {
		return titleStyle.Render("QuickDeck") + "\n\n" + "Terminal too small. Resize to at least 40 columns."
	}
	if model.screen == studyScreen {
		return model.studyView()
	}
	if model.screen == summaryScreen {
		return panelStyle.Render(titleStyle.Render("QuickDeck") + "  " + mutedStyle.Render("Session complete") + "\n\n" + model.deck.Name + "\n\n" + mutedStyle.Render("Enter Return to decks  q Quit"))
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

func (model Model) studyView() string {
	if model.study == nil {
		return panelStyle.Render(titleStyle.Render("QuickDeck") + "\n\nNo study session is active.")
	}
	if model.width > 0 && model.width < 40 {
		return titleStyle.Render("QuickDeck") + "\n\nTerminal too small. Resize to at least 40 columns."
	}
	width := model.width
	if width == 0 {
		width = 80
	}

	card := model.study.Card()
	header := titleStyle.Render("QuickDeck") + "  " + model.deck.Name + "  " + mutedStyle.Render(fmt.Sprintf("Card %d of %d", model.study.Position()+1, model.study.CardCount()))
	answer := "Answer hidden"
	if model.study.IsRevealed() {
		answer = card.Answer
	}
	questionPanel := studyPanel("Question", card.Question, width)
	answerPanel := studyPanel("Answer", answer, width)
	panels := lipgloss.JoinVertical(lipgloss.Left, questionPanel, answerPanel)
	if width >= 100 {
		panelWidth := (width - 2) / 2
		panels = lipgloss.JoinHorizontal(lipgloss.Top,
			studyPanel("Question", card.Question, panelWidth),
			studyPanel("Answer", answer, panelWidth),
		)
	}
	nextAction := "Right/ l Next"
	if model.study.Position() == model.study.CardCount()-1 {
		nextAction = "Right/ l Finish"
	}
	footer := mutedStyle.Render("Left/ h Previous  " + nextAction + "  Space Reveal  q Quit")
	return header + "\n\n" + panels + "\n\n" + footer
}

func studyPanel(title, body string, width int) string {
	innerWidth := width - 6
	if innerWidth < 1 {
		innerWidth = 1
	}
	wrapped := lipgloss.NewStyle().Width(innerWidth).Render(body)
	return panelStyle.Render(titleStyle.Render(title) + "\n\n" + wrapped)
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
