package content

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// LoadManifest reads and validates the deck manifest at the content root.
func LoadManifest(contentDir string) (Manifest, error) {
	manifestPath := filepath.Join(contentDir, "decks.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		if os.IsNotExist(err) {
			return Manifest{}, fmt.Errorf("content dir %q: missing decks.json manifest", contentDir)
		}
		return Manifest{}, fmt.Errorf("content dir %q: unable to read decks.json: %w", contentDir, err)
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return Manifest{}, fmt.Errorf("%s: malformed JSON: %w", manifestPath, err)
	}

	if err := ValidateManifest(manifest, contentDir); err != nil {
		return Manifest{}, fmt.Errorf("%s: %w", manifestPath, err)
	}

	return manifest, nil
}

// LoadCardFile reads and validates a deck card JSON file.
func LoadCardFile(path string) (CardFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return CardFile{}, fmt.Errorf("card file %q: missing file", path)
		}
		return CardFile{}, fmt.Errorf("card file %q: unable to read: %w", path, err)
	}

	var cardFile CardFile
	if err := json.Unmarshal(data, &cardFile); err != nil {
		return CardFile{}, fmt.Errorf("card file %q: malformed JSON: %w", path, err)
	}

	if err := ValidateCardFile(cardFile); err != nil {
		return CardFile{}, fmt.Errorf("card file %q: %w", path, err)
	}

	return cardFile, nil
}
