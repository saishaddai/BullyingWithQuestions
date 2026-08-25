# QuickDeck MVP Specification

**Status:** Draft
**Working name:** QuickDeck
**Audience:** Software engineers revisiting technical concepts
**Primary experience:** A short, offline, keyboard-driven study session

## 1. Product Summary

QuickDeck is a terminal user interface (TUI) for focused study sessions of approximately three minutes. It reads deck metadata and flashcards from local JSON files, lets the user choose a deck, and presents a randomized set of flashcards without repetition.

The application is intentionally small and read-only in the MVP. It does not connect to the internet, use SQL, edit source content, save study progress, or provide social features.

## 2. Goals

- Make it easy to start a short study session from a terminal.
- Help software engineers revisit concepts through question-and-answer flashcards.
- Keep all content local, portable, inspectable, and easy to version-control.
- Provide clear progress and simple keyboard navigation.
- Require no interaction beyond choosing a deck, revealing answers, and moving through cards.

## 3. Non-Goals for MVP

- Internet access or remote content providers.
- SQLite or any other database.
- Creating, editing, deleting, or organizing decks and cards in the app.
- User accounts, profiles, synchronization, social features, or communication.
- Persistent study progress, mastery tracking, scores, or spaced repetition.
- Timers, alarms, or time limits during a session.
- AI-generated cards, AI grading, or external AI services.
- Mouse-driven interaction.
- A web, desktop, or mobile interface.
- Configuration outside the application content directory.

## 4. User Story

As a software engineer, I want to select a topic and review a short randomized set of questions in my terminal so that I can refresh concepts in a few minutes without an internet connection or setup friction.

## 5. Content Model

### Deck

A deck is a freely named technical topic. Decks are displayed as a flat list; grouping and hierarchy are not required.

Required metadata:

- `id`: Stable unique identifier.
- `name`: Human-readable deck name.
- `description`: Brief explanation of the topic.
- `cards_file`: Relative path to the JSON file containing the deck's flashcards.

### Flashcard

A flashcard contains:

- `id`: Stable unique identifier within the deck.
- `question`: Prompt shown before the answer is revealed.
- `answer`: Content shown after the answer is revealed.

The MVP does not require tags, difficulty, notes, sources, multiple answers, or card status. Question and answer values may contain multiple lines. Markdown is treated as plain text unless rendering support is added later.

## 6. JSON Storage Contract

All content is stored under one application content directory. The application does not read content from outside that directory.

Suggested structure:

```text
content/
  decks.json
  cards/
    algorithms.json
    databases.json
    distributed-systems.json
    git.json
    networking.json
    operating-systems.json
    programming-languages.json
```

### `decks.json`

One manifest contains all deck names and descriptions, plus the file reference for each deck's cards.

```json
{
  "version": 1,
  "decks": [
    {
      "id": "algorithms",
      "name": "Algorithms",
      "description": "Core algorithms and complexity concepts.",
      "cards_file": "cards/algorithms.json"
    },
    {
      "id": "databases",
      "name": "Databases",
      "description": "Relational data, indexing, and transactions.",
      "cards_file": "cards/databases.json"
    }
  ]
}
```

The MVP sample content should include at least seven decks. The initial topics may include Algorithms, Databases, Distributed Systems, Git, Networking, Operating Systems, and Programming Languages. The exact content is separate from this product contract and may evolve.

### Deck card file

Each deck has one JSON card file.

```json
{
  "version": 1,
  "deck_id": "algorithms",
  "cards": [
    {
      "id": "algorithms-001",
      "question": "What is the time complexity of binary search on a sorted array?",
      "answer": "O(log n), because the search space is halved after each comparison."
    }
  ]
}
```

Rules:

- `version` is required and must be supported by the application.
- IDs must be non-empty and unique within their collection.
- `cards_file` must be a relative path within the content directory.
- Question and answer must be non-empty strings after trimming whitespace.
- Unknown JSON fields should be ignored for forward compatibility.
- A deck with fewer than 20 valid cards shows all of its valid cards.
- A deck with no valid cards cannot start a study session and must show a useful error.

## 7. Session Behavior

1. On launch, load and validate `decks.json`.
2. Display the available decks in a selectable panel.
3. The user selects a deck with keyboard navigation and confirms the selection.
4. Load and validate the selected deck's card file.
5. Randomly shuffle the valid cards using a session-local random source.
6. Select up to 20 cards from the shuffled list, with no repetition.
7. Show the first question with its answer hidden.
8. The user presses `Space` to reveal the answer for the current card.
9. The user navigates backward and forward through the selected cards.
10. Revealed state belongs to each card in the current session. Revisiting a card keeps its answer revealed.
11. At the end of the session, display a summary.

A deck with 20 or more valid cards always produces exactly 20 session cards. A smaller deck produces all available valid cards. The order is newly randomized each time a session starts.

The application must not repeat a card within a session. Navigating backward to a card does not count as selecting or displaying a new copy of that card.

## 8. TUI Design

The interface should be keyboard-first, readable, and stable while the terminal resizes.

### Deck selection screen

- Header: application name and short status.
- Main panel: flat, scrollable list of decks.
- Detail area: description of the highlighted deck.
- Footer: available keys, such as `Up/Down`, `Enter`, and `q`.
- The selected row must have a clear visual indicator that does not depend on color alone.

### Study screen

Use a two-panel composition:

- Header: selected deck name and session progress, for example `Card 7 of 20`.
- Question panel: the current question, always visible.
- Answer panel: either a placeholder such as `Answer hidden` or the current answer.
- Footer: navigation and reveal state, for example `Left/Right Navigate | Space Reveal | q Quit`.

The question and answer panels should have stable borders and padding. Long content must wrap within its panel and remain readable. The layout must degrade cleanly on narrow terminals; panels may stack vertically when side-by-side rendering is not practical.

### Summary screen

Show:

- Deck name.
- Number of cards reviewed.
- Number of cards revisited by navigating backward.
- Total elapsed session time.
- A clear action to return to deck selection.
- A clear action to quit.

“Cards revisited” means the number of previously viewed cards the user returned to during the session. It is not a count of repeated card copies, because the session contains no duplicate cards.

## 9. Controls

Required MVP controls:

- `Up` / `Down`: Move through the deck list.
- `Enter`: Select a highlighted deck, or confirm a summary action.
- `Left` / `h`: Move to the previous card.
- `Right` / `l`: Move to the next card.
- `Space`: Reveal the current answer.
- `q` or `Ctrl+C`: Quit from any screen.

Navigation should be bounded at the first and last card. The UI must not silently wrap from the last card to the first or vice versa. The user cannot skip cards, rate answers, mark cards, edit content, or start a second branch of interaction inside a session.

Pressing `Space` after the answer is already visible has no additional effect.

## 10. Timing and Summary Metrics

The app does not run a countdown or enforce a three-minute limit. It records the elapsed wall-clock time from the moment the study session begins until the session ends or the summary is displayed.

The summary time should be human-readable, such as `2m 41s`. The app does not persist this value after the process exits.

## 11. Loading and Error Handling

Startup errors must be visible in the TUI or as a clear terminal error before exit. Error messages should identify the affected path and a useful cause.

Required cases:

- Missing content directory: explain where content is expected.
- Missing `decks.json`: explain that the deck manifest is required.
- Malformed JSON: identify the file and report that it could not be parsed.
- Unsupported schema version: identify the file and version.
- Missing card file: identify the deck and referenced path.
- Invalid card records: report the problem. Valid records may still be used when the file can be safely loaded.
- No available decks: explain that at least one valid deck is required.
- No valid cards in a selected deck: prevent the session from starting and return to deck selection.
- Terminal too small: keep controls usable and show a concise resize message when content cannot be rendered safely.

The app should fail closed for unsafe paths and must not follow absolute `cards_file` paths or paths that escape the content directory.

## 12. Operational Constraints

- Offline operation is mandatory; the application must not make network requests.
- Content is read from one local directory only.
- No user data is persisted by the MVP.
- The app should run directly from the Go project and as a compiled executable.
- The existing Bash scripts are legacy utilities and are not part of the QuickDeck MVP workflow. They may be retained temporarily but should not be required for the TUI.

## 13. Acceptance Criteria

- The application launches into a deck-selection panel.
- The manifest can describe at least seven flat, freely named decks.
- A user can select a deck with the keyboard and start a session.
- A deck with 20 or more valid cards displays exactly 20 cards.
- A deck with fewer than 20 valid cards displays all valid cards.
- Cards are randomized on every new session and are not duplicated within that session.
- The question is hidden behind no answer content at the start of each card.
- Pressing `Space` reveals the current answer.
- Left and right navigation moves through the session cards and permits revisiting previous cards.
- The UI shows the current progress, such as `7 of 20`.
- The summary reports elapsed time and the number of cards revisited.
- No database, network connection, editing workflow, or progress persistence is required.
- Missing and malformed content produce actionable errors instead of panics or silent failure.
- Content files outside the configured application directory are rejected.
- The TUI remains usable at common terminal sizes and does not depend on mouse input.

## 14. Suggested Implementation Boundaries

Keep the implementation separable into these responsibilities:

- Content loader: reads the manifest and card files.
- Content validator: validates versions, required fields, IDs, and paths.
- Session builder: shuffles cards and selects the session set.
- Session state: tracks current position, revealed cards, start time, and revisits.
- TUI model and views: renders screens and handles key messages.

The session builder and validator should be testable without starting the TUI. File loading tests should use temporary directories and representative valid and invalid JSON fixtures.

## 15. Future Considerations

These are deliberately deferred and must not expand the MVP:

- Deck and flashcard editing.
- Import and export tools.
- Persistent history and progress.
- Spaced repetition.
- Tags, difficulty, and filtering.
- Configurable content directories.
- Additional interfaces beyond the TUI.
- Optional AI assistance.
- A final public product name and visual identity.

## 16. Open Decisions

- Final product name.
- Exact default content directory name and location.
- Whether the sample seven decks ship in the repository or are supplied separately.
- Whether malformed individual cards should be skipped with a warning or make the whole deck unavailable.
- Whether the summary should count revisits as unique cards revisited or total backward-navigation events.
- Supported Go version and exact Bubble Tea package structure.
