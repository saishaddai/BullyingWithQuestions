package content

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateManifestAcceptsValidDecks(t *testing.T) {
	manifest := Manifest{
		Version: SupportedVersion,
		Decks: []Deck{
			{ID: "algorithms", Name: "Algorithms", Description: "Core algorithms.", CardsFile: "cards/algorithms.json"},
		},
	}

	if err := ValidateManifest(manifest, filepath.Join("content")); err != nil {
		t.Fatalf("ValidateManifest() unexpected error: %v", err)
	}
}

func TestValidateManifestRejectsMissingFields(t *testing.T) {
	manifest := Manifest{
		Version: SupportedVersion,
		Decks: []Deck{{ID: "   ", Name: "Algorithms", Description: "Core algorithms.", CardsFile: "cards/algorithms.json"}},
	}

	if err := ValidateManifest(manifest, filepath.Join("content")); err == nil || !strings.Contains(err.Error(), "deck id") {
		t.Fatalf("ValidateManifest() error = %v, want deck id validation failure", err)
	}
}

func TestValidateManifestRejectsDuplicateDeckIDs(t *testing.T) {
	manifest := Manifest{
		Version: SupportedVersion,
		Decks: []Deck{
			{ID: "algorithms", Name: "Algorithms", Description: "Core algorithms.", CardsFile: "cards/algorithms.json"},
			{ID: "algorithms", Name: "Algorithms 2", Description: "Duplicate.", CardsFile: "cards/algorithms-2.json"},
		},
	}

	if err := ValidateManifest(manifest, filepath.Join("content")); err == nil || !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("ValidateManifest() error = %v, want duplicate deck validation failure", err)
	}
}

func TestValidateManifestRejectsUnsupportedVersion(t *testing.T) {
	manifest := Manifest{Version: 99}
	if err := ValidateManifest(manifest, filepath.Join("content")); err == nil || !strings.Contains(err.Error(), "version") {
		t.Fatalf("ValidateManifest() error = %v, want version validation failure", err)
	}
}

func TestValidateManifestRejectsInvalidCardsFilePaths(t *testing.T) {
	base := filepath.Join("content")
	cases := []struct {
		name     string
		cardsFile string
		want     string
	}{
		{name: "absolute", cardsFile: "/tmp/algorithms.json", want: "relative"},
		{name: "escape", cardsFile: filepath.Join("..", "outside", "algorithms.json"), want: "within"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			manifest := Manifest{
				Version: SupportedVersion,
				Decks: []Deck{{ID: "algorithms", Name: "Algorithms", Description: "Core algorithms.", CardsFile: tc.cardsFile}},
			}

			err := ValidateManifest(manifest, base)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("ValidateManifest(%q) error = %v, want %q", tc.cardsFile, err, tc.want)
			}
		})
	}
}

func TestValidateCardFileAcceptsValidCards(t *testing.T) {
	cardFile := CardFile{
		Version: SupportedVersion,
		DeckID:  "algorithms",
		Cards: []Card{
			{ID: "algorithms-001", Question: "What is binary search?", Answer: "A search that halves the range each step."},
		},
	}

	if err := ValidateCardFile(cardFile); err != nil {
		t.Fatalf("ValidateCardFile() unexpected error: %v", err)
	}
}

func TestValidateCardFileRejectsEmptyValuesAndDuplicates(t *testing.T) {
	cardFile := CardFile{
		Version: SupportedVersion,
		DeckID:  "algorithms",
		Cards: []Card{
			{ID: "algorithms-001", Question: " ", Answer: "A"},
			{ID: "algorithms-001", Question: "A", Answer: "B"},
		},
	}

	if err := ValidateCardFile(cardFile); err == nil || !strings.Contains(err.Error(), "question") && !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("ValidateCardFile() error = %v, want validation failure", err)
	}
}

func TestValidateCardFileRejectsUnsupportedVersion(t *testing.T) {
	cardFile := CardFile{Version: 99, DeckID: "algorithms", Cards: []Card{{ID: "a", Question: "Q", Answer: "A"}}}
	if err := ValidateCardFile(cardFile); err == nil || !strings.Contains(err.Error(), "version") {
		t.Fatalf("ValidateCardFile() error = %v, want version validation failure", err)
	}
}

func TestJSONUnknownFieldsAreIgnored(t *testing.T) {
	data := []byte(`{"version":1,"deck_id":"algorithms","cards":[{"id":"algorithms-001","question":"Q","answer":"A","extra":"ignored"}],"extra":"ignored"}`)

	var cardFile CardFile
	if err := json.Unmarshal(data, &cardFile); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if len(cardFile.Cards) != 1 || cardFile.Cards[0].ID != "algorithms-001" {
		t.Fatalf("json.Unmarshal() did not preserve expected card data: %#v", cardFile)
	}
}

func TestLoadContentDirAcceptsValidDirectory(t *testing.T) {
	dir := t.TempDir()
	cardsDir := filepath.Join(dir, "cards")
	if err := os.MkdirAll(cardsDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	manifest := Manifest{
		Version: SupportedVersion,
		Decks: []Deck{{ID: "algorithms", Name: "Algorithms", Description: "Core algorithms.", CardsFile: "cards/algorithms.json"}},
	}
	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		t.Fatalf("json.Marshal(manifest) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "decks.json"), manifestJSON, 0o644); err != nil {
		t.Fatalf("WriteFile(decks.json) error = %v", err)
	}

	cardFile := CardFile{
		Version: SupportedVersion,
		DeckID:  "algorithms",
		Cards: []Card{{ID: "algorithms-001", Question: "What is binary search?", Answer: "Split the range each step."}},
	}
	cardJSON, err := json.Marshal(cardFile)
	if err != nil {
		t.Fatalf("json.Marshal(cardFile) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(cardsDir, "algorithms.json"), cardJSON, 0o644); err != nil {
		t.Fatalf("WriteFile(card file) error = %v", err)
	}

	loaded, warnings, err := LoadContentDir(dir)
	if err != nil {
		t.Fatalf("LoadContentDir() unexpected error: %v", err)
	}
	if len(loaded) != 1 || len(loaded[0].Cards) != 1 {
		t.Fatalf("LoadContentDir() loaded = %#v, want 1 deck with 1 card", loaded)
	}
	if len(warnings) != 0 {
		t.Fatalf("LoadContentDir() warnings = %v, want none", warnings)
	}
}

func TestLoadContentDirRejectsMissingManifest(t *testing.T) {
	_, _, err := LoadContentDir(t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "decks.json") {
		t.Fatalf("LoadContentDir() error = %v, want decks.json missing error", err)
	}
}

func TestLoadContentDirRejectsEscapingCardPath(t *testing.T) {
	dir := t.TempDir()
	manifest := Manifest{
		Version: SupportedVersion,
		Decks:  []Deck{{ID: "algorithms", Name: "Algorithms", Description: "Core algorithms.", CardsFile: filepath.Join("..", "outside.json")}},
	}
	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		t.Fatalf("json.Marshal(manifest) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "decks.json"), manifestJSON, 0o644); err != nil {
		t.Fatalf("WriteFile(decks.json) error = %v", err)
	}

	_, _, err = LoadContentDir(dir)
	if err == nil || !strings.Contains(err.Error(), "within the content directory") {
		t.Fatalf("LoadContentDir() error = %v, want content-directory safety error", err)
	}
}

func TestLoadContentDirSkipsInvalidCardsWithWarning(t *testing.T) {
	dir := t.TempDir()
	cardsDir := filepath.Join(dir, "cards")
	if err := os.MkdirAll(cardsDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	manifest := Manifest{
		Version: SupportedVersion,
		Decks:  []Deck{{ID: "algorithms", Name: "Algorithms", Description: "Core algorithms.", CardsFile: "cards/algorithms.json"}},
	}
	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		t.Fatalf("json.Marshal(manifest) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "decks.json"), manifestJSON, 0o644); err != nil {
		t.Fatalf("WriteFile(decks.json) error = %v", err)
	}

	cardFile := CardFile{
		Version: SupportedVersion,
		DeckID:  "algorithms",
		Cards: []Card{
			{ID: "algorithms-001", Question: "Valid question", Answer: "Valid answer"},
			{ID: "algorithms-002", Question: " ", Answer: "Missing question"},
			{ID: "algorithms-003", Question: "Another valid question", Answer: "Another valid answer"},
		},
	}
	cardJSON, err := json.Marshal(cardFile)
	if err != nil {
		t.Fatalf("json.Marshal(cardFile) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(cardsDir, "algorithms.json"), cardJSON, 0o644); err != nil {
		t.Fatalf("WriteFile(card file) error = %v", err)
	}

	loaded, warnings, err := LoadContentDir(dir)
	if err != nil {
		t.Fatalf("LoadContentDir() unexpected error: %v", err)
	}
	if len(loaded) != 1 || len(loaded[0].Cards) != 2 {
		t.Fatalf("LoadContentDir() loaded = %#v, want 2 valid cards after filtering", loaded)
	}
	if len(warnings) == 0 {
		t.Fatalf("LoadContentDir() warnings = %v, want at least one warning for skipped invalid cards", warnings)
	}
}
