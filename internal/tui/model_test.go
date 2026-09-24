package tui
package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

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
