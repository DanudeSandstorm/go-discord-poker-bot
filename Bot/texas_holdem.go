package Bot

import (
	"go-poker-bot/Bot/util"
)

// NewTexasHoldem creates a new Texas Hold'em game
func NewTexasHoldem() PokerType {
	suits := []string{Spade, Heart, Diamond, Club}
	ranks := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

	deck := NewDeck(suits, ranks)
	return PokerType{
		GameType: TexasHoldemType,
		Deck:     deck,
		BestHand: TexasHoldemBestHand,
		DealHand: func() []Card {
			return deck.Deal(2)
		},
		String: func() string {
			return "Texas Hold'em"
		},
	}
}

// Returns the best possible 5-card hand that can be made from the five
// community cards and a player's two hole cards
func TexasHoldemBestHand(community []Card, hole []Card) Hand {
	// Combine all cards
	allCards := append(community, hole...)

	// Generate all possible 5-card combinations
	var best Hand
	for handCards := range util.Combinations(allCards, 5) {
		hand := NewHand(handCards)
		if best.Rank == 0 || best.Less(hand) {
			best = hand
		}
	}
	return best
}
