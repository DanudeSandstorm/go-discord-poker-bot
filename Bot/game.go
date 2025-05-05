package Bot

import (
	"github.com/bwmarrin/discordgo"
)

// GameType represents the type of poker game being played
type GameType int

const (
	TexasHoldemType GameType = iota
)

// GameState represents the current state of the game
type GameState int

const (
	NoGame GameState = iota
	Waiting
	NoHands
	HandsDealt
	FlopDealt
	TurnDealt
	RiverDealt
)

// GameOptions represents configurable options for the game
type GameOptions struct {
	SmallBlind int
	BigBlind   int
	MinBuyIn   int
	MaxBuyIn   int
	RaiseDelay int // minutes before blinds double, 0 means off
}

// Game represents the state of a poker game
type Game struct {
	// The type of poker game being played
	Type GameType
	// The pot manager for the game
	PotManager PotManager
	// The deck of cards
	Deck Deck
	// The community cards
	Community []Card
	// The players in the game
	Players []*Player
	// The current state of the game
	State GameState
	// The index of the current dealer
	DealerIndex int
	// The index of the first person to bet in the post-flop rounds
	FirstBettor int
	// The index of the player whose turn it is
	TurnIndex int
	// The players currently in the hand
	InHand []*Player
	// Game options
	Options GameOptions
	// Whether to send all the messages
	Verbose bool
}

// Getters
func (g *Game) GetState() GameState   { return g.State }
func (g *Game) GetPlayers() []*Player { return g.Players }
func (g *Game) GetInHand() []*Player  { return g.InHand }
func (g *Game) GetCommunity() []Card  { return g.Community }
func (g *Game) GetCurrentPlayer() *Player {
	if g.TurnIndex == -1 || g.TurnIndex >= len(g.InHand) {
		return nil
	}
	return g.InHand[g.TurnIndex]
}
func (g *Game) GetOptions() GameOptions { return g.Options }

// Setters
func (g *Game) SetState(state GameState)       { g.State = state }
func (g *Game) SetPlayers(players []*Player)   { g.Players = players }
func (g *Game) SetInHand(inHand []*Player)     { g.InHand = inHand }
func (g *Game) SetCommunity(community []Card)  { g.Community = community }
func (g *Game) SetTurnIndex(index int)         { g.TurnIndex = index }
func (g *Game) SetOptions(options GameOptions) { g.Options = options }

func NewGame() *Game {
	// uses texas holdem by default
	gt := NewTexasHoldem()
	return &Game{
		Type:       TexasHoldemType,
		State:      NoGame,
		PotManager: NewPotManager(gt.BestHand),
		Deck:       gt.Deck,
		Community:  make([]Card, 0),
		Players:    make([]*Player, 0),
		TurnIndex:  -1,
		Options: GameOptions{
			SmallBlind: 1,
			BigBlind:   2,
			MinBuyIn:   50,
			MaxBuyIn:   1000,
			RaiseDelay: 0, // blinds don't raise by default
		},
	}
}

func (g *Game) StartNewGame() {
	g.State = Waiting
	g.Players = make([]*Player, 0)
	g.InHand = make([]*Player, 0)
	g.Community = make([]Card, 0)
	g.TurnIndex = -1
}

func (g Game) IsPlayer(user *discordgo.User) bool {
	for _, p := range g.Players {
		if p.User.ID == user.ID {
			return true
		}
	}
	return false
}

func (g *Game) AddPlayer(user *discordgo.User) {
	if g.IsPlayer(user) {
		return
	}
	g.Players = append(g.Players, &Player{
		User: user,
	})
}

// Stubs for game actions
func (g *Game) DealHands() []string                             { return nil }
func (g *Game) PayBlinds(smallBlind, bigBlind int) []string     { return nil }
func (g *Game) DealFlop() []string                              { return nil }
func (g *Game) DealTurn() []string                              { return nil }
func (g *Game) DealRiver() []string                             { return nil }
func (g *Game) Fold() []string                                  { return nil }
func (g *Game) Call() []string                                  { return nil }
func (g *Game) Raise(amount int) []string                       { return nil }
func (g *Game) Check() []string                                 { return nil }
func (g *Game) Showdown() []string                              { return nil }
func (g *Game) BuyIn(user *discordgo.User, amount int) []string { return nil }
func (g *Game) AllIn() []string                                 { return nil }
func (g *Game) EndGame() []string                               { return nil }
