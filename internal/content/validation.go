package content

import (
	"fmt"
	"path/filepath"
	"strings"
)

const SupportedVersion = 1

// Manifest describes the top-level deck manifest.
type Manifest struct {
	Version int    `json:"version"`
	Decks   []Deck `json:"decks"`
}

// Deck metadata for a single study topic.
type Deck struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CardsFile   string `json:"cards_file"`
}

// CardFile describes a deck's card data payload.
type CardFile struct {
	Version int    `json:"version"`
	DeckID  string `json:"deck_id"`
	Cards   []Card `json:"cards"`
}

// Card is a single flashcard in a deck.
type Card struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// ValidateManifest ensures the manifest matches the required schema.
func ValidateManifest(manifest Manifest, contentDir string) error {
	if manifest.Version != SupportedVersion {
		return fmt.Errorf("manifest version %d is not supported; expected %d", manifest.Version, SupportedVersion)
	}

	seen := make(map[string]struct{}, len(manifest.Decks))
	for i, deck := range manifest.Decks {
		id := strings.TrimSpace(deck.ID)
		if id == "" {
			return fmt.Errorf("manifest deck #%d: deck id is required", i)
		}
		if _, exists := seen[id]; exists {
			return fmt.Errorf("manifest deck %q: duplicate deck id", id)
		}
		seen[id] = struct{}{}

		if strings.TrimSpace(deck.Name) == "" {
			return fmt.Errorf("manifest deck %q: name is required", id)
		}
		if strings.TrimSpace(deck.Description) == "" {
			return fmt.Errorf("manifest deck %q: description is required", id)
		}
		if strings.TrimSpace(deck.CardsFile) == "" {
			return fmt.Errorf("manifest deck %q: cards_file is required", id)
		}

		if err := validateRelativePathWithinContentDir(deck.CardsFile, contentDir, fmt.Sprintf("manifest deck %q", id)); err != nil {
			return err
		}
	}

	return nil
}

// ValidateCardFile ensures a card payload matches the supported schema and card rules.
func ValidateCardFile(cardFile CardFile) error {
	if cardFile.Version != SupportedVersion {
		return fmt.Errorf("card file %q version %d is not supported; expected %d", cardFile.DeckID, cardFile.Version, SupportedVersion)
	}

	deckID := strings.TrimSpace(cardFile.DeckID)
	if deckID == "" {
		return fmt.Errorf("card file deck_id is required")
	}

	seen := make(map[string]struct{}, len(cardFile.Cards))
	for i, card := range cardFile.Cards {
		id := strings.TrimSpace(card.ID)
		if id == "" {
			return fmt.Errorf("card file %q: card #%d id is required", deckID, i)
		}
		if _, exists := seen[id]; exists {
			return fmt.Errorf("card file %q: duplicate card id %q", deckID, id)
		}
		seen[id] = struct{}{}

		if strings.TrimSpace(card.Question) == "" {
			return fmt.Errorf("card file %q: card %q question is required", deckID, id)
		}
		if strings.TrimSpace(card.Answer) == "" {
			return fmt.Errorf("card file %q: card %q answer is required", deckID, id)
		}
	}

	return nil
}

func validateRelativePathWithinContentDir(rawPath, contentDir, context string) error {
	trimmed := strings.TrimSpace(rawPath)
	if trimmed == "" {
		return fmt.Errorf("%s: path is required", context)
	}
	if filepath.IsAbs(trimmed) {
		return fmt.Errorf("%s: cards_file must be relative, got %q", context, rawPath)
	}

	cleaned := filepath.Clean(trimmed)
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) || strings.HasPrefix(cleaned, "..") {
		return fmt.Errorf("%s: cards_file must stay within the content directory: %q", context, rawPath)
	}

	baseAbs, err := filepath.Abs(contentDir)
	if err != nil {
		return fmt.Errorf("%s: unable to resolve content directory %q: %w", context, contentDir, err)
	}
	candidateAbs := filepath.Join(baseAbs, cleaned)
	if rel, err := filepath.Rel(baseAbs, candidateAbs); err == nil && (rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator))) {
		return fmt.Errorf("%s: cards_file must stay within the content directory: %q", context, rawPath)
	}

	return nil
}
