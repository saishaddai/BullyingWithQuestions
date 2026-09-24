package session

import (
	"fmt"
	"math/rand"
	"time"

	"bullyingwithquestions/internal/content"
)

const maxCards = 20

// Session contains the selected cards and state for one study round.
type Session struct {
	cards    []content.Card
	current  int
	revealed map[string]bool
	viewed   map[string]bool
	start    time.Time
	revisits int
}

// New creates a session from a deck's cards. It shuffles a private copy of the
// input and selects at most 20 unique cards.
func New(cards []content.Card, source *rand.Rand, start time.Time) (*Session, error) {
	if len(cards) == 0 {
		return nil, fmt.Errorf("cannot start a session with zero cards")
	}

	selected := append([]content.Card(nil), cards...)
	if source == nil {
		source = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	source.Shuffle(len(selected), func(i, j int) {
		selected[i], selected[j] = selected[j], selected[i]
	})
	if len(selected) > maxCards {
		selected = selected[:maxCards]
	}

	study := &Session{
		cards:    selected,
		revealed: make(map[string]bool, len(selected)),
		viewed:   make(map[string]bool, len(selected)),
		start:    start,
	}
	study.viewed[selected[0].ID] = true
	return study, nil
}

// CardCount returns the number of cards selected for this session.
func (study *Session) CardCount() int {
	return len(study.cards)
}

// Position returns the zero-based position of the current card.
func (study *Session) Position() int {
	return study.current
}

// Card returns the current card.
func (study *Session) Card() content.Card {
	return study.cards[study.current]
}

// CardAt returns the selected card at a zero-based position.
func (study *Session) CardAt(index int) (content.Card, error) {
	if index < 0 || index >= len(study.cards) {
		return content.Card{}, fmt.Errorf("card position %d is out of range", index)
	}
	return study.cards[index], nil
}

func (study *Session) CardAtUnchecked(index int) content.Card {
	return study.cards[index]
}

// IsRevealed reports whether the current card's answer is visible.
func (study *Session) IsRevealed() bool {
	return study.revealed[study.Card().ID]
}

// Reveal reveals the current answer and returns true only when state changes.
func (study *Session) Reveal() bool {
	id := study.Card().ID
	if study.revealed[id] {
		return false
	}
	study.revealed[id] = true
	return true
}

// Next advances to the next card without wrapping.
func (study *Session) Next() bool {
	if study.current >= len(study.cards)-1 {
		return false
	}
	return study.moveTo(study.current + 1)
}

// Previous moves to the previous card without wrapping.
func (study *Session) Previous() bool {
	if study.current <= 0 {
		return false
	}
	return study.moveTo(study.current - 1)
}

func (study *Session) moveTo(position int) bool {
	study.current = position
	id := study.cards[position].ID
	if study.viewed[id] {
		study.revisits++
	} else {
		study.viewed[id] = true
	}
	return true
}

// CardsViewed returns the number of distinct cards reached in this session.
func (study *Session) CardsViewed() int {
	return len(study.viewed)
}

// RevisitCount returns the number of navigations to a previously viewed card.
func (study *Session) RevisitCount() int {
	return study.revisits
}

// StartTime returns when the session began.
func (study *Session) StartTime() time.Time {
	return study.start
}

// ElapsedAt calculates elapsed session time at the supplied instant.
func (study *Session) ElapsedAt(now time.Time) time.Duration {
	if now.Before(study.start) {
		return 0
	}
	return now.Sub(study.start)
}

// Elapsed calculates elapsed session time using the current wall clock.
func (study *Session) Elapsed() time.Duration {
	return study.ElapsedAt(time.Now())
}

// FormatElapsed formats a duration as whole minutes and seconds.
func FormatElapsed(duration time.Duration) string {
	if duration < 0 {
		duration = 0
	}
	seconds := int(duration / time.Second)
	return fmt.Sprintf("%dm %02ds", seconds/60, seconds%60)
}
