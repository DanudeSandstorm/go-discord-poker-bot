package Bot

import (
	"testing"
)

func TestPlayerHandComparison(t *testing.T) {
	tests := []struct {
		name      string
		community []Card
		hole1     []Card
		hole2     []Card
		winner    int
	}{
		{
			name: "High Card vs High Card",
			community: []Card{
				{Suit: Spade, Rank: "9"},
				{Suit: Club, Rank: "4"},
				{Suit: Heart, Rank: "5"},
				{Suit: Spade, Rank: "6"},
				{Suit: Heart, Rank: "7"},
			},
			hole1: []Card{
				{Suit: Club, Rank: "K"},
				{Suit: Club, Rank: "Q"},
			},
			hole2: []Card{
				{Suit: Club, Rank: "A"},
				{Suit: Club, Rank: "2"},
			},
			winner: 2,
		},
		{
			name: "High Card vs High Card (Second Card Decides)",
			community: []Card{
				{Suit: Spade, Rank: "J"},
				{Suit: Club, Rank: "2"},
				{Suit: Heart, Rank: "5"},
				{Suit: Spade, Rank: "9"},
				{Suit: Heart, Rank: "10"},
			},
			hole1: []Card{
				{Suit: Club, Rank: "A"},
				{Suit: Club, Rank: "8"},
			},
			hole2: []Card{
				{Suit: Diamond, Rank: "A"},
				{Suit: Club, Rank: "7"},
			},
			winner: 1,
		},
		{
			name: "Board Beats All Hole Cards",
			community: []Card{
				{Suit: Spade, Rank: "J"},
				{Suit: Club, Rank: "10"},
				{Suit: Heart, Rank: "9"},
				{Suit: Spade, Rank: "8"},
				{Suit: Heart, Rank: "6"},
			},
			hole1: []Card{
				{Suit: Club, Rank: "5"},
				{Suit: Club, Rank: "4"},
			},
			hole2: []Card{
				{Suit: Club, Rank: "3"},
				{Suit: Club, Rank: "2"},
			},
			winner: 0,
		},
		{
			name: "Pair vs High Card",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Club, Rank: "4"},
				{Suit: Heart, Rank: "5"},
				{Suit: Diamond, Rank: "6"},
				{Suit: Club, Rank: "7"},
			},
			hole1: []Card{
				{Suit: Club, Rank: "2"},
				{Suit: Club, Rank: "3"},
			},
			hole2: []Card{
				{Suit: Club, Rank: "A"},
				{Suit: Club, Rank: "K"},
			},
			winner: 1,
		},
		{
			name: "High Pair vs Low Pair",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Diamond, Rank: "3"},
				{Suit: Heart, Rank: "5"},
				{Suit: Diamond, Rank: "6"},
				{Suit: Club, Rank: "7"},
			},
			hole1: []Card{
				{Suit: Club, Rank: "2"},
				{Suit: Club, Rank: "A"},
			},
			hole2: []Card{
				{Suit: Club, Rank: "3"},
				{Suit: Club, Rank: "4"},
			},
			winner: 2,
		},
		{
			name: "Pair with Higher Kicker",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Club, Rank: "4"},
				{Suit: Heart, Rank: "5"},
				{Suit: Diamond, Rank: "6"},
				{Suit: Club, Rank: "7"},
			},
			hole1: []Card{
				{Suit: Club, Rank: "2"},
				{Suit: Club, Rank: "A"},
			},
			hole2: []Card{
				{Suit: Heart, Rank: "2"},
				{Suit: Club, Rank: "K"},
			},
			winner: 1,
		},
		{
			name: "Pair with Low Kickers Tie",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Club, Rank: "K"},
				{Suit: Heart, Rank: "Q"},
				{Suit: Diamond, Rank: "J"},
				{Suit: Club, Rank: "7"},
			},
			hole1: []Card{
				{Suit: Diamond, Rank: "K"},
				{Suit: Club, Rank: "3"},
			},
			hole2: []Card{
				{Suit: Spade, Rank: "K"},
				{Suit: Club, Rank: "6"},
			},
			winner: 0,
		},
		{
			name: "Two Pair vs One Pair",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Club, Rank: "3"},
				{Suit: Heart, Rank: "5"},
				{Suit: Diamond, Rank: "6"},
				{Suit: Club, Rank: "7"},
			},
			hole1: []Card{
				{Suit: Diamond, Rank: "2"},
				{Suit: Spade, Rank: "3"},
			},
			hole2: []Card{
				{Suit: Club, Rank: "A"},
				{Suit: Diamond, Rank: "A"},
			},
			winner: 1,
		},
		{
			name: "Higher Two Pair Wins",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Club, Rank: "4"},
				{Suit: Heart, Rank: "5"},
				{Suit: Diamond, Rank: "6"},
				{Suit: Club, Rank: "7"},
			},
			hole1: []Card{
				{Suit: Club, Rank: "5"},
				{Suit: Spade, Rank: "6"},
			},
			hole2: []Card{
				{Suit: Club, Rank: "2"},
				{Suit: Spade, Rank: "7"},
			},
			winner: 2,
		},
		{
			name: "Two Pair with Higher Second Pair",
			community: []Card{
				{Suit: Club, Rank: "2"},
				{Suit: Spade, Rank: "4"},
				{Suit: Diamond, Rank: "6"},
				{Suit: Heart, Rank: "8"},
				{Suit: Club, Rank: "10"},
			},
			hole1: []Card{
				{Suit: Spade, Rank: "10"},
				{Suit: Club, Rank: "4"},
			},
			hole2: []Card{
				{Suit: Heart, Rank: "10"},
				{Suit: Spade, Rank: "2"},
			},
			winner: 1,
		},
		{
			name: "Two Pair with Higher Kicker",
			community: []Card{
				{Suit: Club, Rank: "K"},
				{Suit: Spade, Rank: "9"},
				{Suit: Heart, Rank: "9"},
				{Suit: Diamond, Rank: "3"},
				{Suit: Spade, Rank: "2"},
			},
			hole1: []Card{
				{Suit: Spade, Rank: "K"},
				{Suit: Club, Rank: "4"},
			},
			hole2: []Card{
				{Suit: Heart, Rank: "K"},
				{Suit: Club, Rank: "5"},
			},
			winner: 2,
		},
		{
			name: "Two Pair with Low Kickers Tie",
			community: []Card{
				{Suit: Club, Rank: "K"},
				{Suit: Spade, Rank: "9"},
				{Suit: Heart, Rank: "9"},
				{Suit: Diamond, Rank: "5"},
				{Suit: Spade, Rank: "2"},
			},
			hole1: []Card{
				{Suit: Spade, Rank: "K"},
				{Suit: Club, Rank: "4"},
			},
			hole2: []Card{
				{Suit: Heart, Rank: "K"},
				{Suit: Club, Rank: "3"},
			},
			winner: 0,
		},
		{
			name: "Three of a Kind vs Pair",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Spade, Rank: "A"},
				{Suit: Club, Rank: "10"},
				{Suit: Heart, Rank: "J"},
				{Suit: Diamond, Rank: "6"},
			},
			hole1: []Card{
				{Suit: Club, Rank: "2"},
				{Suit: Heart, Rank: "2"},
			},
			hole2: []Card{
				{Suit: Diamond, Rank: "A"},
				{Suit: Club, Rank: "K"},
			},
			winner: 1,
		},
		{
			name: "Three of a Kind vs Two Pair",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Spade, Rank: "A"},
				{Suit: Club, Rank: "10"},
				{Suit: Heart, Rank: "K"},
				{Suit: Diamond, Rank: "6"},
			},
			hole1: []Card{
				{Suit: Diamond, Rank: "A"},
				{Suit: Club, Rank: "K"},
			},
			hole2: []Card{
				{Suit: Club, Rank: "2"},
				{Suit: Heart, Rank: "2"},
			},
			winner: 2,
		},
		{
			name: "Higher Three of a Kind",
			community: []Card{
				{Suit: Club, Rank: "K"},
				{Suit: Spade, Rank: "Q"},
				{Suit: Heart, Rank: "5"},
				{Suit: Heart, Rank: "4"},
				{Suit: Heart, Rank: "3"},
			},
			hole1: []Card{
				{Suit: Spade, Rank: "K"},
				{Suit: Heart, Rank: "K"},
			},
			hole2: []Card{
				{Suit: Club, Rank: "Q"},
				{Suit: Heart, Rank: "Q"},
			},
			winner: 1,
		},
		{
			name: "Three of a Kind with Higher Kicker",
			community: []Card{
				{Suit: Club, Rank: "10"},
				{Suit: Spade, Rank: "10"},
				{Suit: Heart, Rank: "J"},
				{Suit: Diamond, Rank: "6"},
				{Suit: Club, Rank: "2"},
			},
			hole1: []Card{
				{Suit: Diamond, Rank: "10"},
				{Suit: Diamond, Rank: "K"},
			},
			hole2: []Card{
				{Suit: Heart, Rank: "10"},
				{Suit: Club, Rank: "A"},
			},
			winner: 2,
		},
		{
			name: "Three of a Kind with Low Kickers Tie",
			community: []Card{
				{Suit: Club, Rank: "10"},
				{Suit: Spade, Rank: "10"},
				{Suit: Heart, Rank: "J"},
				{Suit: Diamond, Rank: "6"},
				{Suit: Club, Rank: "K"},
			},
			hole1: []Card{
				{Suit: Diamond, Rank: "10"},
				{Suit: Diamond, Rank: "3"},
			},
			hole2: []Card{
				{Suit: Heart, Rank: "10"},
				{Suit: Club, Rank: "2"},
			},
			winner: 0,
		},
		{
			name: "Straight vs Pair",
			community: []Card{
				{Suit: Club, Rank: "10"},
				{Suit: Spade, Rank: "9"},
				{Suit: Diamond, Rank: "8"},
				{Suit: Heart, Rank: "7"},
				{Suit: Spade, Rank: "A"},
			},
			hole1: []Card{
				{Suit: Club, Rank: "A"},
				{Suit: Spade, Rank: "K"},
			},
			hole2: []Card{
				{Suit: Club, Rank: "6"},
				{Suit: Club, Rank: "2"},
			},
			winner: 2,
		},
		{
			name: "Straight vs Two Pair",
			community: []Card{
				{Suit: Club, Rank: "10"},
				{Suit: Spade, Rank: "9"},
				{Suit: Diamond, Rank: "8"},
				{Suit: Heart, Rank: "K"},
				{Suit: Spade, Rank: "A"},
			},
			hole1: []Card{
				{Suit: Club, Rank: "6"},
				{Suit: Club, Rank: "7"},
			},
			hole2: []Card{
				{Suit: Club, Rank: "A"},
				{Suit: Spade, Rank: "K"},
			},
			winner: 1,
		},
		{
			name: "Straight vs Three of a Kind",
			community: []Card{
				{Suit: Club, Rank: "10"},
				{Suit: Spade, Rank: "9"},
				{Suit: Diamond, Rank: "8"},
				{Suit: Heart, Rank: "A"},
				{Suit: Spade, Rank: "A"},
			},
			hole1: []Card{
				{Suit: Club, Rank: "A"},
				{Suit: Spade, Rank: "K"},
			},
			hole2: []Card{
				{Suit: Club, Rank: "6"},
				{Suit: Club, Rank: "7"},
			},
			winner: 2,
		},
		{
			name: "Higher Straight",
			community: []Card{
				{Suit: Club, Rank: "10"},
				{Suit: Spade, Rank: "9"},
				{Suit: Diamond, Rank: "8"},
				{Suit: Heart, Rank: "7"},
				{Suit: Spade, Rank: "A"},
			},
			hole1: []Card{
				{Suit: Club, Rank: "J"},
				{Suit: Spade, Rank: "2"},
			},
			hole2: []Card{
				{Suit: Club, Rank: "6"},
				{Suit: Club, Rank: "A"},
			},
			winner: 1,
		},
		{
			name: "Straights Tie",
			community: []Card{
				{Suit: Club, Rank: "10"},
				{Suit: Spade, Rank: "9"},
				{Suit: Diamond, Rank: "8"},
				{Suit: Heart, Rank: "7"},
				{Suit: Spade, Rank: "A"},
			},
			hole1: []Card{
				{Suit: Club, Rank: "6"},
				{Suit: Spade, Rank: "2"},
			},
			hole2: []Card{
				{Suit: Spade, Rank: "6"},
				{Suit: Club, Rank: "A"},
			},
			winner: 0,
		},
		{
			name: "Flush vs Pair",
			community: []Card{
				{Suit: Club, Rank: "5"},
				{Suit: Club, Rank: "6"},
				{Suit: Club, Rank: "7"},
				{Suit: Spade, Rank: "A"},
				{Suit: Spade, Rank: "Q"},
			},
			hole1: []Card{
				{Suit: Club, Rank: "3"},
				{Suit: Club, Rank: "2"},
			},
			hole2: []Card{
				{Suit: Heart, Rank: "A"},
				{Suit: Heart, Rank: "K"},
			},
			winner: 1,
		},
		{
			name: "Flush vs Two Pair",
			community: []Card{
				{Suit: Club, Rank: "5"},
				{Suit: Club, Rank: "6"},
				{Suit: Club, Rank: "7"},
				{Suit: Spade, Rank: "A"},
				{Suit: Spade, Rank: "K"},
			},
			hole1: []Card{
				{Suit: Club, Rank: "3"},
				{Suit: Club, Rank: "2"},
			},
			hole2: []Card{
				{Suit: Heart, Rank: "A"},
				{Suit: Heart, Rank: "K"},
			},
			winner: 1,
		},
		{
			name: "Flush vs Three of a Kind",
			community: []Card{
				{Suit: Club, Rank: "5"},
				{Suit: Club, Rank: "6"},
				{Suit: Club, Rank: "7"},
				{Suit: Spade, Rank: "A"},
				{Suit: Club, Rank: "K"},
			},
			hole1: []Card{
				{Suit: Club, Rank: "3"},
				{Suit: Club, Rank: "2"},
			},
			hole2: []Card{
				{Suit: Heart, Rank: "A"},
				{Suit: Diamond, Rank: "A"},
			},
			winner: 1,
		},
		{
			name: "Flush vs Straight",
			community: []Card{
				{Suit: Club, Rank: "5"},
				{Suit: Club, Rank: "6"},
				{Suit: Club, Rank: "7"},
				{Suit: Spade, Rank: "A"},
				{Suit: Spade, Rank: "Q"},
			},
			hole1: []Card{
				{Suit: Club, Rank: "3"},
				{Suit: Club, Rank: "2"},
			},
			hole2: []Card{
				{Suit: Heart, Rank: "8"},
				{Suit: Heart, Rank: "9"},
			},
			winner: 1,
		},
		{
			name: "Higher Flush",
			community: []Card{
				{Suit: Spade, Rank: "7"},
				{Suit: Spade, Rank: "4"},
				{Suit: Spade, Rank: "8"},
				{Suit: Diamond, Rank: "K"},
				{Suit: Club, Rank: "K"},
			},
			hole1: []Card{
				{Suit: Spade, Rank: "K"},
				{Suit: Spade, Rank: "Q"},
			},
			hole2: []Card{
				{Suit: Spade, Rank: "A"},
				{Suit: Spade, Rank: "2"},
			},
			winner: 2,
		},
		{
			name: "Full House vs Pair",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Club, Rank: "2"},
				{Suit: Heart, Rank: "3"},
				{Suit: Diamond, Rank: "Q"},
				{Suit: Diamond, Rank: "J"},
			},
			hole1: []Card{
				{Suit: Heart, Rank: "2"},
				{Suit: Spade, Rank: "3"},
			},
			hole2: []Card{
				{Suit: Diamond, Rank: "K"},
				{Suit: Diamond, Rank: "A"},
			},
			winner: 1,
		},
		{
			name: "Full House vs Two Pair",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Club, Rank: "2"},
				{Suit: Heart, Rank: "3"},
				{Suit: Diamond, Rank: "Q"},
				{Suit: Diamond, Rank: "A"},
			},
			hole1: []Card{
				{Suit: Diamond, Rank: "K"},
				{Suit: Spade, Rank: "A"},
			},
			hole2: []Card{
				{Suit: Heart, Rank: "2"},
				{Suit: Spade, Rank: "3"},
			},
			winner: 2,
		},
		{
			name: "Full House vs Three of a Kind",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Club, Rank: "2"},
				{Suit: Heart, Rank: "3"},
				{Suit: Diamond, Rank: "Q"},
				{Suit: Diamond, Rank: "A"},
			},
			hole1: []Card{
				{Suit: Heart, Rank: "2"},
				{Suit: Spade, Rank: "3"},
			},
			hole2: []Card{
				{Suit: Diamond, Rank: "K"},
				{Suit: Diamond, Rank: "2"},
			},
			winner: 1,
		},
		{
			name: "Full House vs Straight",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Heart, Rank: "2"},
				{Suit: Club, Rank: "J"},
				{Suit: Diamond, Rank: "10"},
				{Suit: Spade, Rank: "Q"},
			},
			hole1: []Card{
				{Suit: Spade, Rank: "A"},
				{Suit: Spade, Rank: "K"},
			},
			hole2: []Card{
				{Suit: Club, Rank: "2"},
				{Suit: Club, Rank: "Q"},
			},
			winner: 2,
		},
		{
			name: "Full House vs Flush",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Heart, Rank: "2"},
				{Suit: Spade, Rank: "J"},
				{Suit: Diamond, Rank: "10"},
				{Suit: Spade, Rank: "Q"},
			},
			hole1: []Card{
				{Suit: Club, Rank: "2"},
				{Suit: Club, Rank: "Q"},
			},
			hole2: []Card{
				{Suit: Spade, Rank: "A"},
				{Suit: Spade, Rank: "K"},
			},
			winner: 1,
		},
		{
			name: "Full House with Higher Three of a Kind",
			community: []Card{
				{Suit: Spade, Rank: "A"},
				{Suit: Spade, Rank: "K"},
				{Suit: Club, Rank: "J"},
				{Suit: Heart, Rank: "J"},
				{Suit: Diamond, Rank: "3"},
			},
			hole1: []Card{
				{Suit: Club, Rank: "A"},
				{Suit: Diamond, Rank: "J"},
			},
			hole2: []Card{
				{Suit: Club, Rank: "K"},
				{Suit: Diamond, Rank: "K"},
			},
			winner: 2,
		},
		{
			name: "Full House with Higher Pair",
			community: []Card{
				{Suit: Spade, Rank: "A"},
				{Suit: Club, Rank: "A"},
				{Suit: Heart, Rank: "A"},
				{Suit: Club, Rank: "J"},
				{Suit: Club, Rank: "10"},
			},
			hole1: []Card{
				{Suit: Heart, Rank: "J"},
				{Suit: Club, Rank: "2"},
			},
			hole2: []Card{
				{Suit: Heart, Rank: "10"},
				{Suit: Club, Rank: "K"},
			},
			winner: 1,
		},
		{
			name: "Full House Tie",
			community: []Card{
				{Suit: Spade, Rank: "3"},
				{Suit: Club, Rank: "3"},
				{Suit: Heart, Rank: "2"},
				{Suit: Club, Rank: "2"},
				{Suit: Diamond, Rank: "4"},
			},
			hole1: []Card{
				{Suit: Diamond, Rank: "3"},
				{Suit: Spade, Rank: "A"},
			},
			hole2: []Card{
				{Suit: Heart, Rank: "3"},
				{Suit: Club, Rank: "5"},
			},
			winner: 0,
		},
		{
			name: "Four of a Kind vs Pair",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Club, Rank: "2"},
				{Suit: Spade, Rank: "A"},
				{Suit: Heart, Rank: "J"},
				{Suit: Diamond, Rank: "9"},
			},
			hole1: []Card{
				{Suit: Diamond, Rank: "2"},
				{Suit: Heart, Rank: "2"},
			},
			hole2: []Card{
				{Suit: Spade, Rank: "K"},
				{Suit: Spade, Rank: "Q"},
			},
			winner: 1,
		},
		{
			name: "Four of a Kind vs Two Pair",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Club, Rank: "2"},
				{Suit: Spade, Rank: "A"},
				{Suit: Heart, Rank: "J"},
				{Suit: Diamond, Rank: "9"},
			},
			hole1: []Card{
				{Suit: Diamond, Rank: "2"},
				{Suit: Heart, Rank: "2"},
			},
			hole2: []Card{
				{Suit: Spade, Rank: "K"},
				{Suit: Club, Rank: "A"},
			},
			winner: 1,
		},
		{
			name: "Four of a Kind vs Three of a Kind",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Club, Rank: "2"},
				{Suit: Spade, Rank: "Q"},
				{Suit: Heart, Rank: "2"},
				{Suit: Diamond, Rank: "9"},
			},
			hole1: []Card{
				{Suit: Diamond, Rank: "2"},
				{Suit: Heart, Rank: "3"},
			},
			hole2: []Card{
				{Suit: Spade, Rank: "K"},
				{Suit: Club, Rank: "A"},
			},
			winner: 1,
		},
		{
			name: "Four of a Kind vs Straight",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Club, Rank: "2"},
				{Suit: Spade, Rank: "Q"},
				{Suit: Heart, Rank: "J"},
				{Suit: Diamond, Rank: "10"},
			},
			hole1: []Card{
				{Suit: Diamond, Rank: "2"},
				{Suit: Heart, Rank: "2"},
			},
			hole2: []Card{
				{Suit: Spade, Rank: "K"},
				{Suit: Club, Rank: "A"},
			},
			winner: 1,
		},
		{
			name: "Four of a Kind vs Flush",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Club, Rank: "2"},
				{Suit: Spade, Rank: "7"},
				{Suit: Heart, Rank: "J"},
				{Suit: Spade, Rank: "10"},
			},
			hole1: []Card{
				{Suit: Diamond, Rank: "2"},
				{Suit: Heart, Rank: "2"},
			},
			hole2: []Card{
				{Suit: Spade, Rank: "K"},
				{Suit: Spade, Rank: "A"},
			},
			winner: 1,
		},
		{
			name: "Four of a Kind vs Full House",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Club, Rank: "2"},
				{Suit: Spade, Rank: "7"},
				{Suit: Heart, Rank: "J"},
				{Suit: Diamond, Rank: "A"},
			},
			hole1: []Card{
				{Suit: Diamond, Rank: "2"},
				{Suit: Heart, Rank: "2"},
			},
			hole2: []Card{
				{Suit: Diamond, Rank: "7"},
				{Suit: Club, Rank: "7"},
			},
			winner: 1,
		},
		{
			name: "Four of a Kind with Higher Kicker",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Club, Rank: "2"},
				{Suit: Spade, Rank: "7"},
				{Suit: Heart, Rank: "2"},
				{Suit: Diamond, Rank: "2"},
			},
			hole1: []Card{
				{Suit: Diamond, Rank: "K"},
				{Suit: Heart, Rank: "Q"},
			},
			hole2: []Card{
				{Suit: Spade, Rank: "3"},
				{Suit: Spade, Rank: "A"},
			},
			winner: 2,
		},
		{
			name: "Four of a Kind with Same Kicker",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Club, Rank: "2"},
				{Suit: Spade, Rank: "7"},
				{Suit: Heart, Rank: "2"},
				{Suit: Diamond, Rank: "2"},
			},
			hole1: []Card{
				{Suit: Diamond, Rank: "K"},
				{Suit: Heart, Rank: "Q"},
			},
			hole2: []Card{
				{Suit: Spade, Rank: "K"},
				{Suit: Club, Rank: "K"},
			},
			winner: 0,
		},
		{
			name: "Four of a Kind with Low Kickers",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Club, Rank: "2"},
				{Suit: Spade, Rank: "A"},
				{Suit: Heart, Rank: "2"},
				{Suit: Diamond, Rank: "2"},
			},
			hole1: []Card{
				{Suit: Diamond, Rank: "K"},
				{Suit: Heart, Rank: "Q"},
			},
			hole2: []Card{
				{Suit: Spade, Rank: "3"},
				{Suit: Club, Rank: "4"},
			},
			winner: 0,
		},
		{
			name: "Straight Flush vs Pair",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Spade, Rank: "3"},
				{Suit: Spade, Rank: "4"},
				{Suit: Club, Rank: "A"},
				{Suit: Club, Rank: "K"},
			},
			hole1: []Card{
				{Suit: Spade, Rank: "5"},
				{Suit: Spade, Rank: "6"},
			},
			hole2: []Card{
				{Suit: Heart, Rank: "A"},
				{Suit: Heart, Rank: "Q"},
			},
			winner: 1,
		},
		{
			name: "Straight Flush vs Two Pair",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Spade, Rank: "3"},
				{Suit: Spade, Rank: "4"},
				{Suit: Club, Rank: "A"},
				{Suit: Club, Rank: "K"},
			},
			hole1: []Card{
				{Suit: Spade, Rank: "5"},
				{Suit: Spade, Rank: "6"},
			},
			hole2: []Card{
				{Suit: Heart, Rank: "A"},
				{Suit: Heart, Rank: "K"},
			},
			winner: 1,
		},
		{
			name: "Straight Flush vs Three of a Kind",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Spade, Rank: "3"},
				{Suit: Spade, Rank: "4"},
				{Suit: Club, Rank: "A"},
				{Suit: Club, Rank: "K"},
			},
			hole1: []Card{
				{Suit: Spade, Rank: "5"},
				{Suit: Spade, Rank: "6"},
			},
			hole2: []Card{
				{Suit: Heart, Rank: "A"},
				{Suit: Diamond, Rank: "A"},
			},
			winner: 1,
		},
		{
			name: "Straight Flush vs Straight",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Spade, Rank: "3"},
				{Suit: Spade, Rank: "4"},
				{Suit: Club, Rank: "A"},
				{Suit: Club, Rank: "K"},
			},
			hole1: []Card{
				{Suit: Spade, Rank: "5"},
				{Suit: Spade, Rank: "A"},
			},
			hole2: []Card{
				{Suit: Heart, Rank: "5"},
				{Suit: Heart, Rank: "6"},
			},
			winner: 1,
		},
		{
			name: "Straight Flush vs Flush",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Spade, Rank: "3"},
				{Suit: Spade, Rank: "4"},
				{Suit: Club, Rank: "Q"},
				{Suit: Club, Rank: "J"},
			},
			hole1: []Card{
				{Suit: Spade, Rank: "5"},
				{Suit: Spade, Rank: "6"},
			},
			hole2: []Card{
				{Suit: Spade, Rank: "A"},
				{Suit: Spade, Rank: "K"},
			},
			winner: 1,
		},
		{
			name: "Straight Flush vs Full House",
			community: []Card{
				{Suit: Spade, Rank: "2"},
				{Suit: Spade, Rank: "3"},
				{Suit: Spade, Rank: "4"},
				{Suit: Club, Rank: "2"},
				{Suit: Club, Rank: "A"},
			},
			hole1: []Card{
				{Suit: Spade, Rank: "5"},
				{Suit: Spade, Rank: "6"},
			},
			hole2: []Card{
				{Suit: Heart, Rank: "A"},
				{Suit: Diamond, Rank: "A"},
			},
			winner: 1,
		},
		{
			name: "Higher Straight Flush",
			community: []Card{
				{Suit: Spade, Rank: "Q"},
				{Suit: Spade, Rank: "J"},
				{Suit: Spade, Rank: "10"},
				{Suit: Club, Rank: "2"},
				{Suit: Heart, Rank: "3"},
			},
			hole1: []Card{
				{Suit: Spade, Rank: "K"},
				{Suit: Spade, Rank: "A"},
			},
			hole2: []Card{
				{Suit: Spade, Rank: "9"},
				{Suit: Spade, Rank: "8"},
			},
			winner: 1,
		},
		{
			name: "Same Straight Flush Ties",
			community: []Card{
				{Suit: Spade, Rank: "A"},
				{Suit: Spade, Rank: "2"},
				{Suit: Spade, Rank: "3"},
				{Suit: Spade, Rank: "4"},
				{Suit: Spade, Rank: "5"},
			},
			hole1: []Card{
				{Suit: Club, Rank: "6"},
				{Suit: Club, Rank: "7"},
			},
			hole2: []Card{
				{Suit: Heart, Rank: "A"},
				{Suit: Club, Rank: "A"},
			},
			winner: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand1 := TexasHoldemBestHand(tt.community, tt.hole1)
			hand2 := TexasHoldemBestHand(tt.community, tt.hole2)

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
