package Bot

import (
	"math/rand"
	"time"
)

// Deck represents a deck of cards
type Deck struct {
	cards []Card
}

// NewDeck creates a new deck with the given suits and ranks
func NewDeck(suits []string, ranks []string) Deck {
	d := Deck{
		cards: make([]Card, len(suits)*len(ranks)),
	}

	// Initialize all cards
	index := 0
	for _, suit := range suits {
		for _, rank := range ranks {
			d.cards[index] = Card{Suit: suit, Rank: rank}
			index++
		}
	}

	d.Shuffle()

	return d
}

// Shuffle shuffles the deck
func (d *Deck) Shuffle() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

// Deal deals n cards from the top of the deck
func (d *Deck) Deal(n int) []Card {
	if n > len(d.cards) {
		return nil
	}
	cards := d.cards[:n]
	d.cards = d.cards[n:]
	return cards
}
