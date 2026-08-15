package main

type Deck struct {
	id    int
    name  string
    done  bool
}

type Flashcard struct {
    id int
	question string
	answer string
}

type model struct {
    decks         []Deck
    selectedDeck  int
}
