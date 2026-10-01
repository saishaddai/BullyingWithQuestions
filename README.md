# QuickDeck

QuickDeck is an offline, keyboard-driven terminal app for short flashcard study sessions. It reads local JSON decks, shuffles each session, and does not save study progress.

## Requirements

- Go 1.26.5 or newer
- A terminal at least 40 columns wide and 12 rows high for the full interface

## Run

Run commands from the project root so the app can find the `content/` directory:

```sh
go run ./cmd/quickdeck
```

To build and run a standalone executable:

```sh
go build -o quickdeck ./cmd/quickdeck
./quickdeck
```

The executable reads `content/` relative to the current working directory. Keep the project root as the working directory, including when launching a built executable.

## Controls

- `Up` / `Down`: choose a deck
- `Enter`: start a deck; from the summary, return to deck selection
- `Left` / `h`: previous card
- `Right` / `l`: next card; on the last card, finish the session
- `Space`: reveal the current answer
- `q` / `Ctrl+C`: quit

## Content

The `content/` directory contains `decks.json` and the seven sample decks under `content/cards/`. A deck with fewer than 20 valid cards uses all of them; a deck with 20 or more uses 20 shuffled cards. JSON files are read-only inputs. See [SPECS.md](SPECS.md) for the content schema and validation rules.

## Verification

```sh
gofmt -w cmd/quickdeck/*.go internal/app/*.go internal/content/*.go internal/session/*.go internal/tui/*.go
go vet ./...
go test ./...
```

## Project layout

```text
cmd/quickdeck/       executable entry point
content/             sample deck manifest and cards
internal/app/        application lifecycle
internal/content/    JSON loading and validation
internal/session/    randomized session state
internal/tui/         Bubble Tea screens and rendering
```

QuickDeck has no network access, database, content editing, or study-progress persistence.




