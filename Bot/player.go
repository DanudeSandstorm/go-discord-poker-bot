package Bot

import (
	"strings"

	"go-poker-bot/Bot/util"

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

// Returns the amount of money that can be bet by the player
func (p *Player) MaxBet() int {
	return util.Min(p.Balance, p.CurBet+p.Balance)
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
	p.CurBet = util.Min(p.Balance, blind)
	p.Balance -= p.CurBet
	return p.CurBet
}

func (p *Player) PrintHand() string {
	cardsStr := make([]string, len(p.Cards))
	for i, card := range p.Cards {
		cardsStr[i] = card.String()
	}
	return strings.Join(cardsStr, " ")
}
