# BullyingWithQuestions

This project is being rebuilt as a Go-based terminal study app using Bubble Tea. The current Phase 0 baseline sets up a clean application boundary and package layout before the study logic and TUI are implemented.

## Requirements

- Go 1.26.5 or newer
- Standard library only for the current baseline
- Bubble Tea and Lip Gloss will be added in later phases

## Setup

From the project root, run:

```bash
go test ./...
```

To run the application entry point:

```bash
go run ./cmd/quickdeck
```

## Project layout

```text
cmd/quickdeck/      - application entry point
internal/app/       - top-level app lifecycle
internal/content/   - content models and validation
internal/session/    - session state and study rules
internal/tui/       - Bubble Tea model and rendering
```

## Notes

- The MVP is read-only and offline.
- No SQLite or persistence layer is used in the current baseline.
- Phase 0 intentionally leaves the app in a minimal runnable state without deck loading or study behavior.




