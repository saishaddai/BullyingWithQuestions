package app

import (
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"bullyingwithquestions/internal/content"
	"bullyingwithquestions/internal/session"
)

func TestSampleContentBuildsStudySessions(t *testing.T) {
	decks, warnings, err := content.LoadContentDir(filepath.Join("..", "..", "content"))
	if err != nil {
		t.Fatalf("LoadContentDir() error = %v", err)
	}
	if len(decks) < 7 {
		t.Fatalf("loaded %d sample decks, want at least 7", len(decks))
	}
	if len(warnings) != 0 {
		t.Fatalf("sample content warnings = %v, want none", warnings)
	}

	decksByID := make(map[string]content.LoadedDeck, len(decks))
	for _, deck := range decks {
		decksByID[deck.ID] = deck
	}

	largeDeck, ok := decksByID["algorithms"]
	if !ok {
		t.Fatal("sample content is missing the algorithms deck")
	}
	largeSession, err := session.New(largeDeck.Cards, rand.New(rand.NewSource(7)), time.Unix(0, 0))
	if err != nil {
		t.Fatalf("session.New(algorithms) error = %v", err)
	}
	if largeSession.CardCount() != 20 {
		t.Fatalf("large-deck session has %d cards, want 20", largeSession.CardCount())
	}
	selectedIDs := make(map[string]struct{}, largeSession.CardCount())
	for index := 0; index < largeSession.CardCount(); index++ {
		card := largeSession.CardAtUnchecked(index)
		if _, exists := selectedIDs[card.ID]; exists {
			t.Fatalf("large-deck session repeats card %q", card.ID)
		}
		selectedIDs[card.ID] = struct{}{}
	}

	smallDeck, ok := decksByID["git"]
	if !ok {
		t.Fatal("sample content is missing the git deck")
	}
	smallSession, err := session.New(smallDeck.Cards, rand.New(rand.NewSource(8)), time.Unix(0, 0))
	if err != nil {
		t.Fatalf("session.New(git) error = %v", err)
	}
	if smallSession.CardCount() != len(smallDeck.Cards) {
		t.Fatalf("small-deck session has %d cards, want all %d", smallSession.CardCount(), len(smallDeck.Cards))
	}
}

func TestLoadModelBuildsAnErrorModelForMissingContent(t *testing.T) {
	model, err := LoadModel(t.TempDir())
	if err != nil {
		t.Fatalf("LoadModel() unexpected error: %v", err)
	}
	if !strings.Contains(model.View(), "decks.json") {
		t.Fatalf("error view = %q, want decks.json context", model.View())
	}
}

func TestLoadModelExplainsMissingContentDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "missing-content")
	model, err := LoadModel(dir)
	if err != nil {
		t.Fatalf("LoadModel() unexpected error: %v", err)
	}
	view := model.View()
	for _, expected := range []string{"content directory", "does not", "exist;", "decks.json"} {
		if !strings.Contains(view, expected) {
			t.Fatalf("error view omitted %q:\n%s", expected, view)
		}
	}
}

func TestLoadModelPreservesWarningsWhenNoDeckCanBeLoaded(t *testing.T) {
	dir := t.TempDir()
	manifest := `{"version":1,"decks":[{"id":"missing-deck","name":"Missing Deck","description":"Deck with a missing card file.","cards_file":"cards/missing.json"}]}`
	if err := os.WriteFile(filepath.Join(dir, "decks.json"), []byte(manifest), 0o600); err != nil {
		t.Fatalf("WriteFile(decks.json) error = %v", err)
	}

	model, err := LoadModel(dir)
	if err != nil {
		t.Fatalf("LoadModel() unexpected error: %v", err)
	}
	view := model.View()
	for _, expected := range []string{"no valid decks", "available", "missing-deck", "cards/missing.json", "missing file"} {
		if !strings.Contains(view, expected) {
			t.Fatalf("error view omitted %q:\n%s", expected, view)
		}
	}
}
