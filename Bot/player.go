package Bot

import (
	"github.com/bwmarrin/discordgo"
)

type Player struct {
	// How many chips the player has
	Balance int
	// The discord user associated with the player
	User *discordgo.User
	// The player's hole cards
	Cards []Card
	// How many chips the player has bet this round
	CurBet int
	// Whether the player has placed a bet yet this round
	PlacedBet bool
}

// Returns the player's display name
func (p *Player) Name() string {
	if p.User.GlobalName != "" {
		return p.User.GlobalName
	}
	return p.User.Username
}

// Returns the maximum bet that the player can match
func (p *Player) MaxBet() int {
	return p.CurBet + p.Balance
}

// Increases the player's bet to match newAmount
func (p *Player) Bet(newAmount int) int {
	moneyLost := newAmount - p.CurBet
	p.Balance -= moneyLost
	p.CurBet = newAmount
	return moneyLost
}

// Pays the blind amount and returns the amount paid
func (p *Player) PayBlind(blind int) int {
	p.CurBet = min(p.Balance, blind)
	p.Balance -= p.CurBet
	return p.CurBet
}

// Returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
