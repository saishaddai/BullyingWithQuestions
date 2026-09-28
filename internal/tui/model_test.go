package tui

import (
	"math/rand"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"bullyingwithquestions/internal/content"
)

func testDecks(count int) []content.LoadedDeck {
	decks := make([]content.LoadedDeck, count)
	for index := range decks {
		decks[index] = content.LoadedDeck{
			Deck: content.Deck{
				ID:          string(rune('a' + index)),
				Name:        "Deck " + string(rune('A'+index)),
				Description: "A useful study deck.",
			},
			Cards: []content.Card{{ID: "card", Question: "Q", Answer: "A"}},
		}
	}
	return decks
}

func enterStudy(t *testing.T, model Model) Model {
	t.Helper()
	model, command := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if command == nil {
		t.Fatal("Enter returned no deck-selection command")
	}
	model, _ = model.Update(command())
	return model
}

func TestSelectionNavigationIsBounded(t *testing.T) {
	model := NewSelectionModel(testDecks(2), nil)
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	if model.SelectedIndex() != 1 {
		t.Fatalf("selected index = %d, want 1", model.SelectedIndex())
	}
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	if model.SelectedIndex() != 1 {
		t.Fatalf("selection wrapped past last deck: %d", model.SelectedIndex())
	}
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	if model.SelectedIndex() != 0 {
		t.Fatalf("selection moved before first deck: %d", model.SelectedIndex())
	}
}

func TestEnterSelectsDeck(t *testing.T) {
	model := NewSelectionModel(testDecks(2), nil)
	updated, command := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if command == nil {
		t.Fatal("Enter returned no selection command")
	}
	message := command()
	selected, ok := message.(DeckSelectedMsg)
	if !ok || selected.Deck.ID != "a" {
		t.Fatalf("selection message = %#v, want deck a", message)
	}
	if updated.SelectedIndex() != 0 {
		t.Fatal("Enter changed selection unexpectedly")
	}
}

func TestQuitKeyReturnsQuitCommand(t *testing.T) {
	model := NewSelectionModel(testDecks(1), nil)
	_, command := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if command == nil {
		t.Fatal("q returned no quit command")
	}
}

func TestSelectionViewShowsDeckAndDescription(t *testing.T) {
	model := NewSelectionModel(testDecks(1), nil)
	model, _ = model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	view := model.View()
	for _, expected := range []string{"QuickDeck", "Deck A", "A useful study deck.", "[>]", "Up/Down", "Enter", "q"} {
		if !strings.Contains(view, expected) {
			t.Fatalf("View() missing %q in:\n%s", expected, view)
		}
	}
}

func TestNarrowTerminalRendersResizeMessage(t *testing.T) {
	model := NewSelectionModel(testDecks(1), nil)
	model, _ = model.Update(tea.WindowSizeMsg{Width: 20, Height: 10})
	if !strings.Contains(model.View(), "Terminal too small") {
		t.Fatalf("View() did not show narrow-terminal message:\n%s", model.View())
	}
}

func TestErrorModelRendersActionableMessage(t *testing.T) {
	model := NewErrorModel("content/decks.json: missing manifest")
	if !strings.Contains(model.View(), "content/decks.json: missing manifest") {
		t.Fatalf("error view omitted cause:\n%s", model.View())
	}
}

func TestSelectingDeckStartsStudyWithHiddenAnswer(t *testing.T) {
	deck := testDecks(1)[0]
	deck.Cards[0] = content.Card{ID: "card-1", Question: "What is a goroutine?", Answer: "A lightweight concurrent task."}
	model := NewSelectionModel([]content.LoadedDeck{deck}, rand.New(rand.NewSource(1)))
	model, _ = model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model = enterStudy(t, model)

	view := model.View()
	for _, expected := range []string{"Deck A", "Card 1 of 1", "What is a goroutine?", "Answer hidden", "Space Reveal"} {
		if !strings.Contains(view, expected) {
			t.Fatalf("study View() missing %q:\n%s", expected, view)
		}
	}
	if model.screen != studyScreen {
		t.Fatalf("screen = %v, want study", model.screen)
	}
}

func TestStudyRevealAndNavigationKeepPerCardState(t *testing.T) {
	deck := testDecks(1)[0]
	deck.Cards = []content.Card{
		{ID: "card-1", Question: "Question one", Answer: "Answer one"},
		{ID: "card-2", Question: "Question two", Answer: "Answer two"},
	}
	model := enterStudy(t, NewSelectionModel([]content.LoadedDeck{deck}, rand.New(rand.NewSource(1))))
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeySpace})
	firstAnswer := model.study.Card().Answer
	if !model.study.IsRevealed() || !strings.Contains(model.View(), firstAnswer) {
		t.Fatalf("Space did not reveal the current answer:\n%s", model.View())
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	if model.study.Position() != 1 || model.study.IsRevealed() {
		t.Fatal("Right did not advance to the next hidden card")
	}
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	if model.study.Position() != 0 || !model.study.IsRevealed() {
		t.Fatal("h did not return to the revealed first card")
	}
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	if model.study.Position() != 1 {
		t.Fatal("l did not advance to the second card")
	}
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if model.study.Position() != 0 {
		t.Fatal("Left did not navigate to the previous card")
	}
}

func TestAdvancingPastLastCardShowsCompletion(t *testing.T) {
	model := enterStudy(t, NewSelectionModel(testDecks(1), rand.New(rand.NewSource(1))))
	if !strings.Contains(model.View(), "Right/ l Finish") {
		t.Fatalf("final-card footer omits finish action:\n%s", model.View())
	}
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	if model.screen != summaryScreen || !strings.Contains(model.View(), "Session complete") {
		t.Fatalf("last-card advance did not show completion:\n%s", model.View())
	}
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if model.screen != selectionScreen {
		t.Fatalf("Enter from completion returned to screen %v, want selection", model.screen)
	}
}

func TestSummaryShowsMetricsAndRecoveryActions(t *testing.T) {
	deck := testDecks(1)[0]
	deck.Cards = []content.Card{
		{ID: "card-1", Question: "Question one", Answer: "Answer one"},
		{ID: "card-2", Question: "Question two", Answer: "Answer two"},
	}
	model := enterStudy(t, NewSelectionModel([]content.LoadedDeck{deck}, rand.New(rand.NewSource(1))))
	model.study.Next()
	model.study.Previous()
	model.study.Next()
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})

	view := model.View()
	for _, expected := range []string{"Deck A", "Cards reviewed: 2", "Cards revisited: 2", "Elapsed:", "Enter Return to deck selection", "q Quit"} {
		if !strings.Contains(view, expected) {
			t.Fatalf("summary View() missing %q:\n%s", expected, view)
		}
	}
}

func TestEmptySelectedDeckCanReturnToSelection(t *testing.T) {
	deck := testDecks(1)[0]
	deck.Cards = nil
	model := NewSelectionModel([]content.LoadedDeck{deck}, nil)
	model, _ = model.Update(DeckSelectedMsg{Deck: deck})
	if model.screen != errorScreen || !strings.Contains(model.View(), `deck "a"`) || !strings.Contains(model.View(), "Enter Return to deck selection") {
		t.Fatalf("empty-deck error is not actionable:\n%s", model.View())
	}
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if model.screen != selectionScreen {
		t.Fatalf("Enter returned to screen %v, want selection", model.screen)
	}
}

func TestSelectionShowsDetailedContentWarnings(t *testing.T) {
	model := NewSelectionModel(testDecks(1), nil)
	model.SetWarnings([]string{`deck "algorithms": skipped card "bad-card" because question is empty`})
	if view := model.View(); !strings.Contains(view, "bad-card") || !strings.Contains(view, "question is empty") {
		t.Fatalf("selection view omitted warning details:\n%s", view)
	}
}

func TestStudyPanelsWrapWithinTerminalWidth(t *testing.T) {
	deck := testDecks(1)[0]
	deck.Cards[0] = content.Card{
		ID:       "card-1",
		Question: strings.Repeat("A long question with wrapped words. ", 8),
		Answer:   "Answer text",
	}
	model := enterStudy(t, NewSelectionModel([]content.LoadedDeck{deck}, rand.New(rand.NewSource(1))))
	model, _ = model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	for _, line := range strings.Split(model.View(), "\n") {
		if width := lipgloss.Width(line); width > 80 {
			t.Fatalf("study view line width = %d, want at most 80: %q", width, line)
		}
	}
}
