package Bot

const (
	Spade   = "♠"
	Heart   = "♥"
	Diamond = "♦"
	Club    = "♣"
)

type RankInfo struct {
	Name   string
	Plural string
	Value  int
}

var rankInfo = map[string]RankInfo{
	"2":  {"deuce", "deuces", 0},
	"3":  {"three", "threes", 1},
	"4":  {"four", "fours", 2},
	"5":  {"five", "fives", 3},
	"6":  {"six", "sixes", 4},
	"7":  {"seven", "sevens", 5},
	"8":  {"eight", "eights", 6},
	"9":  {"nine", "nines", 7},
	"10": {"ten", "tens", 8},
	"J":  {"jack", "jacks", 9},
	"Q":  {"queen", "queens", 10},
	"K":  {"king", "kings", 11},
	"A":  {"ace", "aces", 12},
}

type Card struct {
	Suit string
	Rank string
}

func (c Card) String() string {
	return c.Suit + c.Rank
}

func (c Card) Name() string {
	return rankInfo[c.Rank].Name
}

func (c Card) Plural() string {
	return rankInfo[c.Rank].Plural
}

func (c Card) Value() int {
	return rankInfo[c.Rank].Value
}

func (c Card) Less(other Card) bool {
	return c.Value() < other.Value()
}

func (c Card) Equal(other Card) bool {
	return c.Rank == other.Rank
}
