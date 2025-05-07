package Bot

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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
	// The last time that the blinds were automatically raised
	LastRaise *time.Time
	// Whether to send all the messages
	Verbose bool
}

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
		LastRaise:  nil,
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
	g.LastRaise = nil
}

func (g *Game) GetState() GameState { return g.State }

func (g *Game) GetPlayers() []*Player { return g.Players }
func (g *Game) GetCurrentPlayer() *Player {
	if g.TurnIndex == -1 || g.TurnIndex >= len(g.InHand) {
		return nil
	}
	return g.InHand[g.TurnIndex]
}

func (g *Game) GetDealer() *Player {
	return g.Players[g.DealerIndex]
}

func (g *Game) GetPlayer(user *discordgo.User) *Player {
	for _, p := range g.Players {
		if p.User.ID == user.ID {
			return p
		}
	}
	return nil
}

func (g Game) IsPlayer(user *discordgo.User) bool {
	for _, p := range g.Players {
		if p.User.ID == user.ID {
			return true
		}
	}
	return false
}

func (g *Game) AddPlayer(user *discordgo.User, name string) {
	g.Players = append(g.Players, &Player{
		User:    user,
		Balance: g.Options.MinBuyIn,
		Name:    name,
	})
}

func (g *Game) BuyIn(user *discordgo.User, amount int, newPlayer bool) []string {
	if amount < g.Options.MinBuyIn {
		return []string{fmt.Sprintf("You must buy in for at least $%d!", g.Options.MinBuyIn)}
	}

	if amount > g.Options.MaxBuyIn {
		return []string{fmt.Sprintf("You can't buy in for more than $%d!", g.Options.MaxBuyIn)}
	}

	// update balance of a player
	player := g.GetPlayer(user)

	if newPlayer {
		return []string{fmt.Sprintf("You've bought in for $%d.", amount)}
	} else {
		return []string{fmt.Sprintf("Increased your balance by $%d. You now have $%d.", amount, player.Balance)}
	}
}

func (g *Game) StatusBetweenRounds() []string {
	messages := []string{}
	if g.Verbose {
		for _, player := range g.Players {
			messages = append(messages, fmt.Sprintf("%s has $%d.", player.Name, player.Balance))
		}
	}
	messages = append(messages, fmt.Sprintf("%s is the current dealer. Message !deal when you're ready.", g.GetDealer().User.Mention()))
	return messages
}

func (g *Game) PayBlinds() []string {
	messages := []string{}

	now := time.Now()
	if g.Options.RaiseDelay == 0 {
		// If the raise delay is set to zero, consider it as being turned
		// off, and do nothing for blinds raises
		g.LastRaise = nil
	} else if g.LastRaise == nil {
		// Start the timer, if it hasn't been started yet
		g.LastRaise = &now
	} else if time.Since(*g.LastRaise) > time.Duration(g.Options.RaiseDelay)*time.Minute {
		messages = append(messages, "**Blinds are being doubled this round!**")
		g.Options.SmallBlind *= 2
		g.Options.BigBlind *= 2
		g.LastRaise = &now
	}

	smallBlind := g.Options.SmallBlind
	bigBlind := g.Options.BigBlind

	var smallPlayer, bigPlayer *Player

	if len(g.InHand) > 2 {
		smallPlayer = g.InHand[(g.DealerIndex+1)%len(g.InHand)]
		bigPlayer = g.InHand[(g.DealerIndex+2)%len(g.InHand)]
		// The first player to bet pre-flop is the player to the left of the big blind
		g.TurnIndex = (g.DealerIndex + 3) % len(g.InHand)
		// The first player to bet post-flop is the first player to the left of the dealer
		g.FirstBettor = (g.DealerIndex + 1) % len(g.Players)
	} else {
		// In heads-up games, who plays the blinds is different, with the
		// dealer playing the small blind and the other player paying the big
		smallPlayer = g.Players[g.DealerIndex]
		bigPlayer = g.Players[(g.DealerIndex+1)%len(g.Players)]

		// Dealer goes first pre-flop, the other player goes first afterwards
		g.TurnIndex = g.DealerIndex
		g.FirstBettor = (g.DealerIndex + 1) % len(g.Players)
	}

	messages = append(messages, fmt.Sprintf("%s has paid the small blind of $%d.", smallPlayer.Name, smallBlind))

	if g.PotManager.PayBlind(smallPlayer, smallBlind) {
		messages = append(messages, fmt.Sprintf("%s is all in!", smallPlayer.Name))
		g.LeaveHand(smallPlayer)
	}

	messages = append(messages, fmt.Sprintf("%s has paid the big blind of $%d.", bigPlayer.Name, bigBlind))

	if g.PotManager.PayBlind(bigPlayer, bigBlind) {
		messages = append(messages, fmt.Sprintf("%s is all in!", bigPlayer.Name))
		g.LeaveHand(bigPlayer)
	}

	return messages
}

func (g *Game) Showdown() []string {
	messages := []string{}

	for len(g.Community) < 5 {
		g.Community = append(g.Community, g.Deck.Deal(1)...)
	}

	messages = append(messages, "We have reached the end of betting. "+
		"All cards will be revealed.")

	communityStr := make([]string, len(g.Community))
	for i, card := range g.Community {
		communityStr[i] = card.String()
	}
	messages = append(messages, strings.Join(communityStr, "  "))

	for player := range g.PotManager.InPot() {
		messages = append(messages, fmt.Sprintf("%s's hand: %s", player.Name, player.PrintHand()))
	}

	winners := g.PotManager.GetWinners(g.Community)

	for winner, winnings := range winners {
		handName := g.PotManager.BestHandFunc(winner.Cards, g.Community)
		messages = append(messages, fmt.Sprintf("%s wins $%d with a %s.", winner.Name, winnings, handName))
		winner.Balance += winnings
	}

	// Remove players that went all in and lost
	i := 0
	for i < len(g.Players) {
		player := g.Players[i]
		if player.Balance > 0 {
			i++
		} else {
			messages = append(messages, fmt.Sprintf("%s has been knocked out of the game!", player.Name))
			g.Players = append(g.Players[:i], g.Players[i+1:]...)
			if len(g.Players) == 1 {
				// There's only one player, so they win
				messages = append(messages, fmt.Sprintf("%s wins the game! Congratulations!", g.Players[0].Name))
				g.State = NoGame
				return messages
			}
			if i <= g.DealerIndex {
				g.DealerIndex -= 1
			}
		}
	}

	// Go on to the next round
	g.State = NoHands
	g.NextDealer()
	return append(messages, g.StatusBetweenRounds()...)
}

func (g *Game) NextDealer() {
	g.DealerIndex = (g.DealerIndex + 1) % len(g.Players)
}

func (g *Game) Fold() []string {
	var messages []string

	if g.Verbose {
		messages = append(messages, fmt.Sprintf("%s has folded.", g.GetCurrentPlayer().Name))
	}

	g.PotManager.HandleFold(g.GetCurrentPlayer())
	g.LeaveHand(g.GetCurrentPlayer())

	// If only one person is left in the pot, give it to them instantly
	if len(g.PotManager.InPot()) == 1 {
		// grab the only player in the pot
		var winner *Player
		for p := range g.PotManager.InPot() {
			winner = p
			break
		}
		messages = append(messages, fmt.Sprintf("%s wins $%d!", winner.Name, g.PotManager.Value()))
		winner.Balance += g.PotManager.Value()
		g.State = NoHands
		g.NextDealer()
		return append(messages, g.StatusBetweenRounds()...)
	}

	// If there's still betting to do, go on to the next turn
	if !g.PotManager.BettingOver() {
		g.TurnIndex -= 1
		return append(messages, g.NextTurn()...)
	}

	// Otherwise, have the showdown immediately
	return g.Showdown()
}

func (g *Game) Call() []string {
	messages := []string{}

	g.PotManager.HandleCall(g.GetCurrentPlayer())

	if g.Verbose {
		messages = append(messages, fmt.Sprintf("%s calls.", g.GetCurrentPlayer().Name))
	}

	if g.GetCurrentPlayer().Balance == 0 {
		messages = append(messages, fmt.Sprintf("%s is all in!", g.GetCurrentPlayer().Name))
		g.LeaveHand(g.GetCurrentPlayer())
		g.TurnIndex -= 1
	}

	return append(messages, g.NextTurn()...)
}

func (g *Game) Raise(amount int) []string {
	messages := []string{}

	g.PotManager.HandleRaise(g.GetCurrentPlayer(), amount)

	if g.Verbose {
		messages = append(messages, fmt.Sprintf("%s raises by $%d.", g.GetCurrentPlayer().Name, amount))
	}

	if g.GetCurrentPlayer().Balance == 0 {
		messages = append(messages, fmt.Sprintf("%s is all in!", g.GetCurrentPlayer().Name))
		g.LeaveHand(g.GetCurrentPlayer())
		g.TurnIndex -= 1
	}

	return append(messages, g.NextTurn()...)
}

func (g *Game) Check() []string {
	messages := []string{}

	g.GetCurrentPlayer().PlacedBet = true

	if g.Verbose {
		messages = append(messages, fmt.Sprintf("%s calls.", g.GetCurrentPlayer().Name))
	}

	return append(messages, g.NextTurn()...)
}

func (g *Game) AllIn() []string {
	if g.PotManager.CurBet() > g.GetCurrentPlayer().MaxBet() {
		return g.Call()
	} else {
		return g.Raise(g.GetCurrentPlayer().MaxBet() - g.PotManager.CurBet())
	}
}

// Removes a player from being able to bet, if they folded or went all in
func (g *Game) LeaveHand(player *Player) {
	for i, p := range g.InHand {
		if p == player {
			g.InHand = append(g.InHand[:i], g.InHand[i+1:]...)

			// Adjust the index of the first person to bet and the index of the
			// current player, depending on the index of the player who just folded
			if i < g.FirstBettor {
				g.FirstBettor -= 1
			}
			if g.FirstBettor >= len(g.InHand) {
				g.FirstBettor = 0
			}
			if g.TurnIndex >= len(g.InHand) {
				g.TurnIndex = 0
			}
			return
		}
	}
}

func (g *Game) NextRound() []string {
	messages := []string{}

	switch g.State {
	case HandsDealt:
		messages = append(messages, "Dealing the flop:")
		g.Community = append(g.Community, g.Deck.Deal(3)...)
		g.State = FlopDealt
	case FlopDealt:
		messages = append(messages, "Dealing the turn:")
		g.Community = append(g.Community, g.Deck.Deal(1)...)
		g.State = TurnDealt
	case TurnDealt:
		messages = append(messages, "Dealing the river:")
		g.Community = append(g.Community, g.Deck.Deal(1)...)
		g.State = RiverDealt
	case RiverDealt:
		return g.Showdown()
	}

	communityStr := make([]string, len(g.Community))
	for i, card := range g.Community {
		communityStr[i] = card.String()
	}
	messages = append(messages, strings.Join(communityStr, "  "))

	g.PotManager.NextRound()
	g.TurnIndex = g.FirstBettor

	return append(messages, g.CurOptions()...)
}

func (g *Game) NextTurn() []string {
	if g.PotManager.RoundOver() {
		if g.PotManager.BettingOver() {
			return g.Showdown()
		} else {
			return g.NextRound()
		}
	}

	g.TurnIndex = (g.TurnIndex + 1) % len(g.InHand)
	return g.CurOptions()
}

func (g *Game) CurOptions() []string {
	messages := []string{
		fmt.Sprintf("It is %s's turn. Current balance is $%d.",
			g.GetCurrentPlayer().User.Mention(),
			g.GetCurrentPlayer().Balance,
		),
	}

	curBet := g.PotManager.CurBet()
	if curBet > 0 {
		messages = append(messages, fmt.Sprintf("The pot is currently $%d. The current bet to meet is $%d, and %s has bet $%d.",
			g.PotManager.Value(),
			curBet,
			g.GetCurrentPlayer().Name,
			g.GetCurrentPlayer().CurBet))
	} else {
		messages = append(messages, fmt.Sprintf("The pot is currently $%d. The current bet to meet is $%d.",
			g.PotManager.Value(),
			curBet,
		))
	}

	if g.Verbose {
		if g.GetCurrentPlayer().CurBet == curBet {
			messages = append(messages, "Message !check, !raise or !fold.")
		} else if g.GetCurrentPlayer().MaxBet() > curBet {
			messages = append(messages, "Message !call, !raise or !fold.")
		} else {
			messages = append(messages, "Message !allin or !fold.")
		}
	}

	return messages
}

func (g *Game) DealHands() []string {
	g.Deck.Shuffle()

	// Start out the shared cards as being empty
	g.Community = make([]Card, 0)

	// Deals hands to each player, setting their initial bets to zero and
	// adding them as being in on the hand
	g.InHand = make([]*Player, 0)
	for _, player := range g.Players {
		player.Cards = g.Deck.Deal(2)
		player.CurBet = 0
		player.PlacedBet = false
		g.InHand = append(g.InHand, player)
	}

	g.State = HandsDealt
	messages := []string{"The hands have been dealt!"}

	// Reset the pot for the new hand
	g.PotManager.NewHand(g.Players)

	// Pay blinds if there are any
	if g.Options.SmallBlind > 0 {
		messages = append(messages, g.PayBlinds()...)
	}
	g.TurnIndex--
	messages = append(messages, g.NextTurn()...)
	return messages
}

// EndGame ends the current game and returns messages about final chip counts
func (g *Game) EndGame() []string {
	messages := []string{"Game has been ended."}
	for _, player := range g.Players {
		messages = append(messages, fmt.Sprintf("%s has $%d.", player.Name, player.Balance))
	}

	g.State = NoGame
	return messages
}

// ListOptions returns a string listing the current game options
func (g *Game) ListOptions() string {
	return fmt.Sprintf("Current game options:\n"+
		"Small Blind: $%d\n"+
		"Big Blind: $%d\n"+
		"Min Buy-In: $%d\n"+
		"Max Buy-In: $%d\n"+
		"Blind Raise Delay: %d minutes (0 = off)",
		g.Options.SmallBlind, g.Options.BigBlind, g.Options.MinBuyIn, g.Options.MaxBuyIn, g.Options.RaiseDelay)
}

// HandleOptions handles the options command and returns messages to be sent
func (g *Game) SetOption(args []string) string {
	option := strings.ToLower(args[0])
	amount, err := strconv.Atoi(args[1])
	if err != nil {
		return "Invalid amount!"
	}

	switch option {
	case "sb":
		if amount <= 0 {
			return "Small blind must be greater than 0!"
		}
		if amount >= g.Options.BigBlind {
			return "Small blind must be less than big blind!"
		}
		g.Options.SmallBlind = amount
	case "bb":
		if amount <= 0 {
			return "Small blind must be greater than 0!"
		}
		if amount < g.Options.SmallBlind {
			return "Big blind must be greater than or equal to the small blind!"
		}
		g.Options.BigBlind = amount
	case "min":
		if amount <= 0 {
			return "Min buy-in must be greater than 0!"
		}
		if amount >= g.Options.MaxBuyIn {
			return "Min buy-in must be less than max buy-in!"
		}
		g.Options.MinBuyIn = amount
	case "max":
		if amount <= g.Options.MinBuyIn {
			return "Max buy-in must be greater than min buy-in!"
		}
		g.Options.MaxBuyIn = amount
	case "delay":
		if amount < 0 {
			return "Delay must be 0 or greater!"
		}
		g.Options.RaiseDelay = amount
	default:
		return "Invalid option! Use sb, bb, min, max, or delay"
	}

	return fmt.Sprintf("%s set to %d", option, amount)
}

func (g *Game) ToggleVerbose() string {
	g.Verbose = !g.Verbose
	return fmt.Sprintf("Verbose mode is now %t", g.Verbose)
}

func (g *Game) IsCurrentPlayer(user *discordgo.User) bool {
	currentPlayer := g.GetCurrentPlayer()
	return currentPlayer != nil && currentPlayer.User.ID == user.ID
}
