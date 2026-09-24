package app

import (
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
