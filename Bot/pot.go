package Bot

import (
	"go-poker-bot/Bot/util"
	"math"
)

type Pot struct {
	// The players that have contributed to this pot and can win it
	Players map[*Player]struct{}
	// The bet that needs to be made to join this pot
	CurBet int
	// The amount of money accumulated in this pot
	Amount int
	// The maximum bet that can be held by this pot before it needs a side pot
	MaxBet int
}

func NewPot(players map[*Player]struct{}) Pot {
	p := Pot{
		Players: players,
		CurBet:  0,
		Amount:  0,
	}

	if len(players) > 0 {
		// Find the minimum max bet among players
		minMaxBet := math.MaxInt
		for player := range players {
			if player.MaxBet() < minMaxBet {
				minMaxBet = player.MaxBet()
			}
		}
		p.MaxBet = minMaxBet
	} else {
		// Set an extremely high max bet to prevent accidental creation of side pot when blinds are impossibly high
		p.MaxBet = math.MaxInt
	}

	return p
}

// Returns which players win this pot, based on the given community cards
func (p Pot) GetWinners(community []Card, bestHandFunc BestHandFunc) []*Player {
	var winners []*Player
	var bestHand Hand

	for player := range p.Players {
		hand := bestHandFunc(community, player.Cards)
		if bestHand.Rank == 0 || bestHand.Less(hand) {
			winners = []*Player{player}
			bestHand = hand
		} else if hand.Equal(bestHand) {
			winners = append(winners, player)
		}
	}
	return winners
}

// Returns a new side pot, for when the bet overflows what can be contained
// in this pot
func (p Pot) MakeSidePot() Pot {
	excluded := make(map[*Player]struct{})
	for player := range p.Players {
		if player.MaxBet() == p.MaxBet {
			excluded[player] = struct{}{}
		}
	}

	// Create new players set without excluded players
	newPlayers := make(map[*Player]struct{})
	for player := range p.Players {
		if _, ok := excluded[player]; !ok {
			newPlayers[player] = struct{}{}
		}
	}

	return NewPot(newPlayers)
}

type PotManager struct {
	// List of side pots in the game
	// If nobody's all-in, there should only be one pot
	// Higher-priced pots are towards the end of the list
	Pots    []Pot
	LastBet int
}

func NewPotManager() PotManager {
	return PotManager{
		Pots: make([]Pot, 0),
	}
}

// Resets the list of pots for a new hand
func (pm *PotManager) NewHand(players []*Player) {
	// Convert slice to map for Pot creation
	playerSet := make(map[*Player]struct{})
	for _, player := range players {
		playerSet[player] = struct{}{}
	}
	pm.Pots = []Pot{NewPot(playerSet)}
}

// Returns the current bet to be matched
func (pm PotManager) CurBet() int {
	total := 0
	for _, pot := range pm.Pots {
		total += pot.CurBet
	}
	return total
}

// Returns the amount of money that's in all the pots and side pots
func (pm PotManager) Value() int {
	total := 0
	for _, pot := range pm.Pots {
		total += pot.Amount
	}
	return total
}

// Increases the current bet to a new given amount
func (pm *PotManager) IncreaseBet(newAmount int) {
	accumulatedBet := 0
	for pm.Pots[len(pm.Pots)-1].MaxBet < newAmount {
		pm.Pots[len(pm.Pots)-1].CurBet = pm.Pots[len(pm.Pots)-1].MaxBet - accumulatedBet
		accumulatedBet += pm.Pots[len(pm.Pots)-1].CurBet
		pm.Pots = append(pm.Pots, pm.Pots[len(pm.Pots)-1].MakeSidePot())
	}
	newBet := util.Min(pm.Pots[len(pm.Pots)-1].MaxBet, newAmount)
	pm.Pots[len(pm.Pots)-1].CurBet = newBet - accumulatedBet
	pm.LastBet = newAmount
}

// Returns all the players that are in the pot
func (pm PotManager) InPot() map[*Player]struct{} {
	return pm.Pots[0].Players
}

// Handles a player folding, removing them from every pot they're eligible for
func (pm *PotManager) HandleFold(player *Player) {
	for i := range pm.Pots {
		delete(pm.Pots[i].Players, player)
	}
}

// Handles a player calling the current bet
func (pm *PotManager) HandleCall(player *Player) {
	newAmount := player.Bet(util.Min(player.MaxBet(), pm.CurBet()))
	oldBet := player.CurBet - newAmount
	potIndex := 0
	for newAmount > 0 {
		curPot := &pm.Pots[potIndex]
		oldBet -= curPot.CurBet
		if oldBet < 0 {
			curPot.Amount -= oldBet
			newAmount += oldBet
			oldBet = 0
		}
		potIndex++
	}
	player.PlacedBet = true
}

// Handles a player raising the current bet to a given amount
func (pm *PotManager) HandleRaise(player *Player, newAmount int) {
	pm.IncreaseBet(pm.CurBet() + newAmount)
	pm.HandleCall(player)
}

// Pays the initial blinds for the player, returning whether they were
// forced to go all-in by the blind
func (pm *PotManager) PayBlind(player *Player, blind int) bool {
	pm.IncreaseBet(blind)
	pm.HandleCall(player)
	player.PlacedBet = false
	return player.Balance == 0
}

// Returns whether the betting round is over
func (pm PotManager) RoundOver() bool {
	if pm.BettingOver() {
		pm.LastBet = 0
		return true
	}
	for player := range pm.Pots[0].Players {
		if player.Balance == 0 {
			continue
		}
		if !player.PlacedBet || player.CurBet < pm.CurBet() {
			return false
		}
	}
	pm.LastBet = 0
	return true
}

// Returns whether all betting is over
func (pm PotManager) BettingOver() bool {
	playersLeftBetting := false
	for player := range pm.Pots[0].Players {
		if player.Balance > 0 {
			if playersLeftBetting || !player.PlacedBet {
				return false
			}
			if player.CurBet < pm.CurBet() {
				return false
			}
			playersLeftBetting = true
		}
	}
	return true
}

// Returns the winners of the pot, and the amounts that they won
func (pm PotManager) GetWinners(sharedCards []Card, bestHandFunc BestHandFunc) map[*Player]int {
	winners := make(map[*Player]int)
	for _, pot := range pm.Pots {
		potWinners := pot.GetWinners(sharedCards, bestHandFunc)
		if len(potWinners) == 0 {
			continue
		}
		potWon := pot.Amount / len(potWinners)
		if potWon > 0 {
			for _, winner := range potWinners {
				winners[winner] += potWon
			}
		}
	}
	return winners
}

// Advances to the next round of betting
func (pm *PotManager) NextRound() {
	for i := range pm.Pots {
		pm.Pots[i].CurBet = 0
		pm.Pots[i].MaxBet = 0
	}
	for player := range pm.Pots[len(pm.Pots)-1].Players {
		player.PlacedBet = false
		player.CurBet = 0
	}

	// Find the minimum max bet among remaining players
	minMaxBet := math.MaxInt
	for player := range pm.Pots[len(pm.Pots)-1].Players {
		if player.MaxBet() < minMaxBet {
			minMaxBet = player.MaxBet()
		}
	}
	pm.Pots[len(pm.Pots)-1].MaxBet = minMaxBet
}
