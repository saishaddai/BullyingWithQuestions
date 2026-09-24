package session

import (
	"math/rand"
	"testing"
	"time"

	"bullyingwithquestions/internal/content"
)

func cards(count int) []content.Card {
	result := make([]content.Card, count)
	for index := range result {
		result[index] = content.Card{
			ID:       string(rune('a' + index)),
			Question: "Question",
			Answer:   "Answer",
		}
	}
	return result
}

func TestNewRejectsEmptyDeck(t *testing.T) {
	if _, err := New(nil, rand.New(rand.NewSource(1)), time.Unix(100, 0)); err == nil {
		t.Fatal("New() accepted an empty deck")
	}
}

func TestNewKeepsSmallDeckCards(t *testing.T) {
	input := cards(3)
	study, err := New(input, rand.New(rand.NewSource(1)), time.Unix(100, 0))
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
	if study.CardCount() != len(input) {
		t.Fatalf("CardCount() = %d, want %d", study.CardCount(), len(input))
	}
}

func TestNewLimitsLargeDeckToTwentyUniqueCards(t *testing.T) {
	study, err := New(cards(25), rand.New(rand.NewSource(1)), time.Unix(100, 0))
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
	if study.CardCount() != 20 {
		t.Fatalf("CardCount() = %d, want 20", study.CardCount())
	}

	seen := make(map[string]bool)
	for index := 0; index < study.CardCount(); index++ {
		card, err := study.CardAt(index)
		if err != nil {
			t.Fatalf("CardAt(%d) unexpected error: %v", index, err)
		}
		if seen[card.ID] {
			t.Fatalf("duplicate selected card %q", card.ID)
		}
		seen[card.ID] = true
	}
}

func TestNewDoesNotMutateInputAndRandomizesUsingLocalSource(t *testing.T) {
	input := cards(25)
	originalFirst := input[0].ID
	first, err := New(input, rand.New(rand.NewSource(1)), time.Unix(100, 0))
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
	second, err := New(input, rand.New(rand.NewSource(2)), time.Unix(100, 0))
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
	if input[0].ID != originalFirst {
		t.Fatal("New() mutated the input card order")
	}
	if first.CardAtUnchecked(0).ID == second.CardAtUnchecked(0).ID {
		t.Fatal("different session-local random sources produced the same first card")
	}
}

func TestRevealStateSurvivesBoundedNavigation(t *testing.T) {
	study, err := New(cards(2), rand.New(rand.NewSource(1)), time.Unix(100, 0))
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	if study.IsRevealed() {
		t.Fatal("new card is already revealed")
	}
	if !study.Reveal() || study.Reveal() {
		t.Fatal("Reveal() should change state once and then be idempotent")
	}
	if !study.Next() || study.Next() {
		t.Fatal("Next() should be bounded at the last card")
	}
	if !study.Previous() || study.Previous() {
		t.Fatal("Previous() should be bounded at the first card")
	}
	if !study.IsRevealed() {
		t.Fatal("reveal state was lost when navigating back")
	}
}

func TestNavigationTracksViewedAndRevisitedCards(t *testing.T) {
	study, err := New(cards(3), rand.New(rand.NewSource(1)), time.Unix(100, 0))
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
	if study.CardsViewed() != 1 || study.RevisitCount() != 0 {
		t.Fatalf("initial metrics = viewed %d, revisits %d", study.CardsViewed(), study.RevisitCount())
	}

	study.Next()
	study.Next()
	if study.CardsViewed() != 3 || study.RevisitCount() != 0 {
		t.Fatalf("forward metrics = viewed %d, revisits %d", study.CardsViewed(), study.RevisitCount())
	}
	study.Previous()
	if study.RevisitCount() != 1 {
		t.Fatalf("RevisitCount() = %d after returning to a viewed card, want 1", study.RevisitCount())
	}
}

func TestElapsedTimeFormatting(t *testing.T) {
	study, err := New(cards(1), rand.New(rand.NewSource(1)), time.Unix(100, 0))
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
	if got := study.ElapsedAt(time.Unix(221, 500_000_000)); got != 2*time.Minute+1*time.Second+500*time.Millisecond {
		t.Fatalf("ElapsedAt() = %v, want 2m1.5s", got)
	}
	if got := FormatElapsed(2*time.Minute + 41*time.Second); got != "2m 41s" {
		t.Fatalf("FormatElapsed() = %q, want %q", got, "2m 41s")
	}
}
