package Bot

import (
	"go-poker-bot/Bot/util"
)

// Returns the best possible 5-card hand that can be made from the five
// community cards and a player's two hole cards
func bestPossibleHand(community []Card, hole [2]Card) Hand {
	// Combine all cards
	allCards := append(community, hole[:]...)

	// Generate all possible 5-card combinations
	var best Hand
	for handCards := range util.Combinations(allCards, 5) {
		hand := NewHand(handCards)
		if best.Rank == 0 || hand.Less(best) {
			best = hand
		}
	}
	return best
}
