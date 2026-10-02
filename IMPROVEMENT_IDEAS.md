# QuickDeck Improvement Ideas

These proposals expand beyond the current offline flashcard MVP while keeping every workflow inside the terminal UI. They are independent directions; each can be scoped and delivered separately.

## 1. Persistent Progress and Spaced Review

Let learners rate each answer (for example: again, difficult, good, or easy) and use those ratings to schedule future reviews. Add a TUI dashboard for cards due today, upcoming reviews, and per-deck progress.

A first milestone could save review dates and ratings to a local versioned JSON file, then add a due-card session mode alongside the current shuffled session. Keep the existing no-history session available, make saved data inspectable and exportable, and write files atomically so an interrupted save does not corrupt progress.

**Why it expands the project:** The current app forgets each session; this adds a continuing learning loop, user data, scheduling rules, and persistence while remaining offline.

## 2. In-TUI Deck and Card Studio

Add keyboard-driven workflows to create, edit, duplicate, and remove decks and flashcards without leaving QuickDeck. Provide a form-like editor, field validation, a preview of the resulting card, and a review screen before writing changes.

Start with editing a card and adding a card to an existing deck. Preserve the current JSON format, protect against accidental overwrite, and use atomic writes. Later, add import/export of deck files and a validation-only preview for content brought in from elsewhere.

**Why it expands the project:** Content is currently read-only and must be edited outside the application; a TUI studio turns QuickDeck into a local content-management tool without introducing another interface.

## 3. Searchable Library and Custom Session Builder

Make large collections easier to navigate with incremental deck/card search, tags, and filters such as topic or difficulty. Let learners build a session from matching cards and choose its size or selection strategy instead of always using the current fixed limit of 20.

A first milestone could add deck search and configurable session length, followed by card-level tags and a filter preview that shows how many cards match before the session starts. Keep selection usable with the keyboard and provide clear empty-result states.

**Why it expands the project:** The current library is a flat deck list and sessions use a fixed randomized selection; this supports larger, more targeted study libraries while keeping the existing quick-start path.
