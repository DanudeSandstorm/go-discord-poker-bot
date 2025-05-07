package Bot

import (
	"testing"
)

func TestOmahaHandComparison(t *testing.T) {
	tests := []struct {
		name      string
		community []Card
		hole1     []Card
		hole2     []Card
		winner    int
	}{
		{
			name: "Full House Beats Ace High Straight",
			community: []Card{
				{Suit: Spade, Rank: "A"},
				{Suit: Heart, Rank: "A"},
				{Suit: Diamond, Rank: "A"},
				{Suit: Club, Rank: "K"},
				{Suit: Spade, Rank: "Q"},
			},
			hole1: []Card{
				{Suit: Spade, Rank: "J"},
				{Suit: Heart, Rank: "10"},
				{Suit: Diamond, Rank: "9"},
				{Suit: Club, Rank: "8"},
			},
			hole2: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Heart, Rank: "3"},
				{Suit: Diamond, Rank: "5"},
				{Suit: Club, Rank: "5"},
			},
			winner: 2,
		},
		{
			name: "Straight Flush beats Flush",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Spade, Rank: "3"},
				{Suit: Spade, Rank: "4"},
				{Suit: Club, Rank: "A"},
				{Suit: Heart, Rank: "K"},
			},
			hole1: []Card{
				{Suit: Spade, Rank: "5"},
				{Suit: Spade, Rank: "6"},
				{Suit: Heart, Rank: "7"},
				{Suit: Heart, Rank: "8"},
			},
			hole2: []Card{
				{Suit: Spade, Rank: "7"},
				{Suit: Spade, Rank: "8"},
				{Suit: Heart, Rank: "9"},
				{Suit: Heart, Rank: "10"},
			},
			winner: 1,
		},
		{
			name: "Straight Beats Trip Aces",
			community: []Card{
				{Suit: Spade, Rank: "A"},
				{Suit: Heart, Rank: "A"},
				{Suit: Diamond, Rank: "K"},
				{Suit: Club, Rank: "K"},
				{Suit: Spade, Rank: "Q"},
			},
			hole1: []Card{
				{Suit: Spade, Rank: "J"},
				{Suit: Heart, Rank: "10"},
				{Suit: Diamond, Rank: "9"},
				{Suit: Club, Rank: "8"},
			},
			hole2: []Card{
				{Suit: Spade, Rank: "7"},
				{Suit: Heart, Rank: "6"},
				{Suit: Diamond, Rank: "A"},
				{Suit: Club, Rank: "2"},
			},
			winner: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand1 := OmahaBestHand(tt.community, tt.hole1)
			hand2 := OmahaBestHand(tt.community, tt.hole2)

			if tt.winner == 1 {
				if hand1.Less(hand2) {
					t.Errorf("player1 should win but hand1.Less(hand2) returned true")
					t.Logf("player1 hand: %v", hand1)
					t.Logf("player2 hand: %v", hand2)
				}
			} else if tt.winner == 2 {
				if !hand1.Less(hand2) {
					t.Errorf("player2 should win but hand1.Less(hand2) returned false")
					t.Logf("player1 hand: %v", hand1)
					t.Logf("player2 hand: %v", hand2)
				}
			} else {
				// should be a tie
				if hand1.Less(hand2) {
					t.Errorf("should be a tie but hand1.Less(hand2) returned true")
					t.Logf("player1 hand: %v", hand1)
					t.Logf("player2 hand: %v", hand2)
				} else if hand2.Less(hand1) {
					t.Errorf("should be a tie but hand2.Less(hand1) returned true")
					t.Logf("player1 hand: %v", hand1)
					t.Logf("player2 hand: %v", hand2)
				}
			}
		})
	}
}
