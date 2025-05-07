package Bot

import (
	"go-poker-bot/Bot/util"
)

// NewPotLimitOmaha creates a new Pot Limit Omaha game
func NewPotLimitOmaha() PokerType {
	suits := []string{Spade, Heart, Diamond, Club}
	ranks := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

	deck := NewDeck(suits, ranks)
	return PokerType{
		GameType: PotLimitOmahaType,
		Deck:     deck,
		BestHand: OmahaBestHand,
		DealHand: func() []Card {
			return deck.Deal(4)
		},
	}
}

// Returns the best possible 5-card hand that can be made from the five
// community cards and a player's four hole cards, using exactly 2 hole cards
// and exactly 3 community cards
func OmahaBestHand(community []Card, hole []Card) Hand {
	if len(hole) != 4 {
		panic("Omaha requires exactly 4 hole cards")
	}
	if len(community) != 5 {
		panic("Omaha requires exactly 5 community cards")
	}

	// Generate all possible combinations of 2 hole cards
	var best Hand
	for holeCards := range util.Combinations(hole, 2) {
		// For each combination of 2 hole cards, try all combinations of 3 community cards
		for commCards := range util.Combinations(community, 3) {
			// Combine the selected cards
			handCards := append(holeCards, commCards...)
			hand := NewHand(handCards)

			// Update best hand if this is better
			if best.Less(hand) {
				best = hand
			}
		}
	}
	return best
}
