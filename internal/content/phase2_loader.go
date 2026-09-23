package content

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadedDeck represents a manifest deck with validated cards loaded from disk.
type LoadedDeck struct {
	Deck
	Cards []Card
}

// LoadContentDir reads a content directory and returns all valid decks plus warnings
// about skipped invalid records.
func LoadContentDir(contentDir string) ([]LoadedDeck, []string, error) {
	manifest, err := LoadManifest(contentDir)
	if err != nil {
		return nil, nil, err
	}

	loaded := make([]LoadedDeck, 0, len(manifest.Decks))
	warnings := make([]string, 0)

	for _, deck := range manifest.Decks {
		cardPath := filepath.Join(contentDir, deck.CardsFile)
		validCards, deckWarnings, err := loadDeckCards(cardPath, deck.ID)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("deck %q (%s): %v", deck.ID, deck.CardsFile, err))
			continue
		}
		warnings = append(warnings, deckWarnings...)
		if len(validCards) == 0 {
			warnings = append(warnings, fmt.Sprintf("deck %q: no valid cards available", deck.ID))
			continue
		}

		loaded = append(loaded, LoadedDeck{Deck: deck, Cards: validCards})
	}

	if len(loaded) == 0 {
		return nil, warnings, fmt.Errorf("content dir %q: no valid decks available", contentDir)
	}

	return loaded, warnings, nil
}

// LoadSampleContent loads the default content directory used by the app.
func LoadSampleContent() ([]LoadedDeck, []string, error) {
	return LoadContentDir("content")
}

func loadDeckCards(cardPath, deckID string) ([]Card, []string, error) {
	data, err := os.ReadFile(cardPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, fmt.Errorf("card file %q: missing file", cardPath)
		}
		return nil, nil, fmt.Errorf("card file %q: unable to read: %w", cardPath, err)
	}

	var raw struct {
		Version int            `json:"version"`
		DeckID  string         `json:"deck_id"`
		Cards   []json.RawMessage `json:"cards"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, nil, fmt.Errorf("card file %q: malformed JSON: %w", cardPath, err)
	}
	if raw.Version != SupportedVersion {
		return nil, nil, fmt.Errorf("card file %q version %d is not supported; expected %d", deckID, raw.Version, SupportedVersion)
	}
	if strings.TrimSpace(raw.DeckID) != "" && strings.TrimSpace(raw.DeckID) != strings.TrimSpace(deckID) {
		return nil, nil, fmt.Errorf("card file %q: deck_id %q does not match manifest deck id %q", cardPath, raw.DeckID, deckID)
	}

	validCards := make([]Card, 0, len(raw.Cards))
	seenIDs := make(map[string]struct{}, len(raw.Cards))
	warnings := make([]string, 0)
	for i, rawCard := range raw.Cards {
		var card Card
		if err := json.Unmarshal(rawCard, &card); err != nil {
			warnings = append(warnings, fmt.Sprintf("deck %q: skipped malformed card #%d: %v", deckID, i, err))
			continue
		}
		id := strings.TrimSpace(card.ID)
		if id == "" {
			warnings = append(warnings, fmt.Sprintf("deck %q: skipped card #%d because id is empty", deckID, i))
			continue
		}
		if strings.TrimSpace(card.Question) == "" {
			warnings = append(warnings, fmt.Sprintf("deck %q: skipped card %q because question is empty", deckID, id))
			continue
		}
		if strings.TrimSpace(card.Answer) == "" {
			warnings = append(warnings, fmt.Sprintf("deck %q: skipped card %q because answer is empty", deckID, id))
			continue
		}
		if _, exists := seenIDs[id]; exists {
			warnings = append(warnings, fmt.Sprintf("deck %q: skipped duplicate card id %q", deckID, id))
			continue
		}
		seenIDs[id] = struct{}{}
		validCards = append(validCards, card)
	}

	return validCards, warnings, nil
}
