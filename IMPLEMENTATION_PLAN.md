# QuickDeck Implementation Plan

This plan turns `SPECS.md` into incremental phases for the Go implementation. Each phase should leave the project in a runnable or testable state and should be completed before the next phase begins.

## Guiding Decisions

- Use Go's standard library for JSON parsing, filesystem access, validation, and randomness.
- Use Bubble Tea for the TUI and Lip Gloss for layout and styling.
- Keep content loading, validation, session behavior, and rendering separated so the non-TUI behavior is easy to test.
- Treat the MVP as read-only and offline. Do not introduce SQLite, network clients, persistence, editing, or scoring.
- Preserve the existing Bash utilities as legacy files, but do not make them part of the QuickDeck runtime.

## Phase 0: Project Baseline

**Goal:** Establish a clean Go application boundary before adding behavior.

### Work

- Confirm the supported Go version and update the module name if `menu` is not the intended module path.
- Replace the placeholder structs in `main.go` with a minimal application entry point.
- Organize packages around the suggested responsibilities, for example:
  - `internal/content`: models, loading, validation, and safe path handling.
  - `internal/session`: shuffled card selection and session state.
  - `internal/tui`: Bubble Tea model, messages, views, and key handling.
- Add a small `cmd/quickdeck` entry point if the project benefits from separating executable code from packages.
- Add formatting and test commands to the README.

### Exit Check

`go test ./...` and `go vet ./...` pass with no study functionality implemented yet.

## Phase 1: Content Model and Validation

**Goal:** Represent the JSON contract and reject invalid data deterministically.

### Work

- Define exported or package-level types for:
  - Manifest and deck metadata.
  - Card files and flashcards.
  - Supported schema version.
- Validate required fields after trimming whitespace.
- Validate non-empty and unique IDs within manifests and card collections.
- Validate supported `version` values.
- Validate that `cards_file` is relative and remains inside the configured content directory.
- Ignore unknown JSON fields for forward compatibility.
- Return errors that include the affected file or deck and the specific cause.

### Tests

- Valid manifest and card file.
- Missing required fields.
- Duplicate IDs.
- Unsupported versions.
- Empty question or answer values.
- Absolute paths and paths that escape the content directory.
- Unknown fields.

### Exit Check

Validator tests pass without starting Bubble Tea.

## Phase 2: File Loading and Sample Content

**Goal:** Load a complete local content directory safely.

### Work

- Choose and document the default content directory, such as `content/` beside the executable's working directory.
- Load `decks.json` and report missing or malformed files clearly.
- Load each referenced card file only after manifest validation.
- Allow valid card records to remain usable when invalid records can be skipped safely, while surfacing a useful warning.
- Reject an empty or entirely invalid deck for study.
- Add the seven initial sample decks described in the specification:
  - Algorithms
  - Databases
  - Distributed Systems
  - Git
  - Networking
  - Operating Systems
  - Programming Languages
- Include enough cards to exercise both the fewer-than-20 and at-least-20 cases.

### Tests

- Use temporary directories for missing, malformed, incomplete, and valid content fixtures.
- Verify the loaded deck count and card count.
- Verify that referenced files outside the content directory are rejected.

### Exit Check

The loader can read the sample content and produces actionable errors for the required failure cases.

## Phase 3: Session Builder and State

**Goal:** Implement the study rules independently of the UI.

### Work

- Shuffle cards with a session-local random source.
- Select exactly 20 cards when a deck has 20 or more valid cards; otherwise select every valid card.
- Guarantee no duplicate card IDs in a session.
- Track:
  - Current card position.
  - Revealed state per selected card.
  - Session start time.
  - Cards viewed.
  - Revisit count according to the final interpretation in `SPECS.md`.
- Make first/last navigation bounded rather than wrapping.
- Make reveal idempotent.
- Calculate elapsed time and format it as human-readable minutes and seconds.

### Tests

- Zero-card deck is rejected.
- Small deck returns all cards.
- Large deck returns exactly 20 unique cards.
- Repeated sessions can produce different orders without relying on global random state.
- Reveal state survives backward and forward navigation.
- Navigation does not wrap.
- Revisit metrics and elapsed-time formatting are correct.

### Exit Check

All session behavior is covered by deterministic unit tests without terminal rendering.

## Phase 4: Deck Selection Screen

**Goal:** Make the first user workflow usable in the terminal.

### Work

- Create the Bubble Tea application model with explicit screens: selection, study, summary, and error.
- Render the application header and short status.
- Render a flat, scrollable deck list with a non-color selection indicator.
- Render the highlighted deck description.
- Support `Up`, `Down`, `Enter`, and `q`/`Ctrl+C`.
- Handle an empty deck list and deck-loading errors without panicking.
- Respond to terminal resize messages and preserve stable layout dimensions.

### Exit Check

Launching the program shows the deck-selection screen, and a user can select a valid deck or quit.

## Phase 5: Study Screen and Navigation

**Goal:** Connect the tested session state to the required study experience.

### Work

- Start a new session after selecting a deck.
- Render the selected deck name and progress such as `Card 7 of 20`.
- Render question and answer panels with stable borders, padding, and wrapped multiline content.
- Show `Answer hidden` until `Space` is pressed.
- Support `Left`/`h`, `Right`/`l`, `Space`, and quit controls.
- Stack panels or show a concise resize message when the terminal is too narrow.
- Ensure the final card transitions to the summary screen through an explicit user action or the defined end-of-session behavior.

### Exit Check

A complete session can be started, revealed, navigated in both directions, and completed without losing per-card reveal state.

## Phase 6: Summary and Error UX

**Goal:** Finish the user-facing MVP workflows and failure paths.

### Work

- Render the selected deck name, reviewed-card count, revisit count, and elapsed time.
- Provide clear summary actions to return to deck selection or quit.
- Make startup errors visible before exit or in the TUI, including:
  - Missing content directory.
  - Missing or malformed manifest.
  - Unsupported versions.
  - Missing card files.
  - Invalid card records.
  - No available decks.
  - No valid cards in a selected deck.
  - Terminal too small.
- Keep error messages tied to paths and deck IDs where applicable.

### Exit Check

Every required error case from `SPECS.md` has an actionable message and a safe recovery or exit path.

## Phase 7: Integration, Hardening, and Release Readiness

**Goal:** Verify the full MVP and make it straightforward to run.

### Work

- Add integration tests for loading sample content and constructing a session.
- Test common terminal sizes manually or with a terminal harness where practical.
- Run `gofmt`, `go vet ./...`, and `go test ./...`.
- Run the application from the Go project and from a compiled binary.
- Update `README.md` with the actual setup, content directory, controls, and commands.
- Remove outdated SQLite requirements from the README because SQLite is outside the MVP.
- Confirm there are no network calls, database dependencies, progress files, or content writes.
- Review the implementation against every acceptance criterion in `SPECS.md`.

### Exit Check

The application launches into deck selection, supports a full offline study session, handles malformed content safely, and passes the documented verification commands.

## Suggested Delivery Order

1. Phase 0: Project baseline.
2. Phase 1: Models and validation.
3. Phase 2: Loading and sample content.
4. Phase 3: Session logic.
5. Phase 4: Deck selection UI.
6. Phase 5: Study UI.
7. Phase 6: Summary and error UX.
8. Phase 7: Integration and release readiness.

This order keeps the highest-risk rules, especially validation, safe paths, random selection, navigation, and revisit counting, testable before they are embedded in terminal rendering.

## Definition of Done

- The program starts in a keyboard-driven deck-selection screen.
- The sample content includes at least seven valid decks.
- Sessions contain up to 20 unique cards and preserve reveal state while navigating.
- Progress, revisits, elapsed time, and summary actions are visible.
- Invalid or unsafe content produces useful errors without panics.
- The application is offline, read-only, database-free, and does not persist study data.
- The README describes the actual Go and TUI workflow.