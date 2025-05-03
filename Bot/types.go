package Bot

// BestHandFunc defines the signature for functions that determine the best possible hand
// given community cards and hole cards
type BestHandFunc func(community []Card, hole []Card) Hand

type PokerType struct {
	Deck     Deck
	BestHand BestHandFunc
}
