package Bot

import (
	"slices"
	"sort"
)

type HandRanking int

const (
	HighCard HandRanking = iota + 1
	Pair
	TwoPair
	ThreeOfKind
	Straight
	Flush
	FullHouse
	FourOfKind
	StraightFlush
)

type Hand struct {
	Cards []Card
	Rank  HandRanking
}

func NewHand(cards []Card) Hand {
	// Sort the cards to make comparison easier
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Less(cards[j])
	})

	h := Hand{
		Cards: cards,
	}

	// Get duplicates (pairs, three-of-a-kinds, etc)
	dups := h.getDuplicates()

	// Determine hand ranking
	if h.isFlush() {
		if h.isStraight() {
			h.Rank = StraightFlush
		} else {
			h.Rank = Flush
		}
	} else if h.isStraight() {
		h.Rank = Straight
	} else if len(dups) > 0 {
		if len(dups) == 2 {
			if len(dups[1]) == 3 {
				h.Rank = FullHouse
			} else {
				h.Rank = TwoPair
			}
		} else {
			if len(dups[0]) == 4 {
				h.Rank = FourOfKind
			} else if len(dups[0]) == 3 {
				h.Rank = ThreeOfKind
			} else {
				h.Rank = Pair
			}
		}
		h.rearrangeDuplicates(dups)
	} else {
		h.Rank = HighCard
	}

	return h
}

func (h Hand) String() string {
	switch h.Rank {
	case HighCard:
		return h.Cards[4].Name() + " high"
	case Pair:
		return "pair of " + h.Cards[4].Plural()
	case TwoPair:
		return "two pair, " + h.Cards[4].Plural() + " and " + h.Cards[2].Plural()
	case ThreeOfKind:
		return "three of a kind, " + h.Cards[4].Plural()
	case Straight:
		return h.Cards[4].Name() + "-high straight"
	case Flush:
		return h.Cards[4].Name() + "-high flush"
	case FullHouse:
		return "full house, " + h.Cards[4].Plural() + " over " + h.Cards[1].Plural()
	case FourOfKind:
		return "four of a kind, " + h.Cards[4].Plural()
	case StraightFlush:
		if h.Cards[4].Rank == "A" {
			return "royal flush"
		}
		return h.Cards[4].Name() + "-high straight flush"
	default:
		return "unknown hand"
	}
}

func (h Hand) Less(other Hand) bool {
	if h.Rank < other.Rank {
		return true
	}
	if h.Rank > other.Rank {
		return false
	}
	// Compare cards from highest to lowest
	for i := len(h.Cards) - 1; i >= 0; i-- {
		if h.Cards[i].Less(other.Cards[i]) {
			return true
		}
		if other.Cards[i].Less(h.Cards[i]) {
			return false
		}
	}
	return false
}

func (h Hand) Equal(other Hand) bool {
	if h.Rank != other.Rank {
		return false
	}
	for i := range h.Cards {
		if !h.Cards[i].Equal(other.Cards[i]) {
			return false
		}
	}
	return true
}

func (h Hand) isStraight() bool {
	// Check normal straight
	for i := 1; i < 5; i++ {
		if h.Cards[i-1].Value() != h.Cards[i].Value()-1 {
			break
		}
		if i == 4 {
			return true
		}
	}

	// Check ace-low straight
	values := make([]int, 5)
	for i, card := range h.Cards {
		values[i] = card.Value()
	}
	return values[0] == 0 && values[1] == 1 && values[2] == 2 && values[3] == 3 && values[4] == 12
}

func (h Hand) isFlush() bool {
	suit := h.Cards[0].Suit
	for _, card := range h.Cards[1:] {
		if card.Suit != suit {
			return false
		}
	}
	return true
}

// Returns a list of the pairs, three-of-a-kinds and four-of-a-kinds in the hand
func (h Hand) getDuplicates() [][]Card {
	var dups [][]Card
	var curDup []Card

	for i, card := range h.Cards {
		if i == 0 {
			curDup = []Card{card}
			continue
		}

		if !card.Equal(curDup[0]) {
			if len(curDup) > 1 {
				dups = append(dups, curDup)
			}
			curDup = []Card{card}
		} else {
			curDup = append(curDup, card)
		}
	}

	if len(curDup) > 1 {
		dups = append(dups, curDup)
	}

	// For calculating full houses, ensure three-of-a-kind is second
	if len(dups) == 2 && len(dups[0]) > len(dups[1]) {
		dups[0], dups[1] = dups[1], dups[0]
	}

	return dups
}

// Rearrange the duplicated cards in the hand so that comparing two hands
// with the same ranking is easier
// This moves duplicated cards to the end of the hand
func (h *Hand) rearrangeDuplicates(dups [][]Card) {
	var flatDups []Card
	for _, dup := range dups {
		flatDups = append(flatDups, dup...)
	}

	for _, dup := range flatDups {
		for i, card := range h.Cards {
			if card.Equal(dup) {
				h.Cards = slices.Delete(h.Cards, i, i+1)
				break
			}
		}
	}

	h.Cards = append(h.Cards, flatDups...)
}
