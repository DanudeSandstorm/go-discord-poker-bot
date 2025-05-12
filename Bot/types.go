package Bot

// GameType represents the type of poker game being played
type GameType int

const (
	TexasHoldemType GameType = iota
	PotLimitOmahaType
)

// BestHandFunc defines the signature for functions that determine the best possible hand
// given community cards and hole cards
type BestHandFunc func(community []Card, hole []Card) Hand

type PokerType struct {
	GameType GameType
	Deck     Deck
	BestHand BestHandFunc
	DealHand func() []Card
	String   func() string
	MaxBet   func(player *Player, pm *PotManager) int
}
