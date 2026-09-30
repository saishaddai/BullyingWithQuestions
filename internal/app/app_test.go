package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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
