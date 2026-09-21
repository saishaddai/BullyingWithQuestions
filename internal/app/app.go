package app
package app

import (
	"fmt"
	"os"
)

// Run provides the application boundary for the QuickDeck MVP and returns nil until
// the actual TUI lifecycle is implemented.
func Run() error {
	// Phase 0 deliberately leaves the app executable but does not start the TUI yet.
	// The future UI and deck-loading behavior will live in this package boundary.
	if _, err := fmt.Fprintln(os.Stdout, "QuickDeck baseline initialized"); err != nil {
		return err
	}
	return nil
}
