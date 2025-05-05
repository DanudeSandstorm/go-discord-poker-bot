package Bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	games map[string]*Game
}

func NewBot() *Bot {
	return &Bot{
		games: make(map[string]*Game),
	}
}

func checkNilErr(e error) {
	if e != nil {
		log.Fatal("Error message")
	}
}

func Run(token string) {
	discord, err := discordgo.New("Bot " + token)
	checkNilErr(err)

	bot := NewBot()
	discord.AddHandler(bot.newMessage)

	discord.Open()
	defer discord.Close()

	fmt.Println("Bot running....")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
}

func (b *Bot) getGame(channelID string) *Game {
	game, exists := b.games[channelID]
	if !exists {
		game = NewGame()
		b.games[channelID] = game
	}
	return game
}

func (b *Bot) newMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	// Ignore messages that don't start with !
	if !strings.HasPrefix(m.Content, "!") {
		return
	}

	// Split message into command and arguments
	parts := strings.Fields(m.Content[1:]) // Strip the ! from the first part
	// TODO allow custom prefixes
	if len(parts) == 0 {
		return
	}

	command := strings.ToLower(parts[0])
	args := parts[1:]

	// Map shorthand commands to their full versions
	shorthand := map[string]string{
		"b":   "raise",
		"bet": "raise",
		"c":   "call",
		"d":   "deal",
		"f":   "fold",
		"r":   "raise",
		"x":   "check",
	}

	// If command is a shorthand, replace it with the full version
	if fullCmd, exists := shorthand[command]; exists {
		command = fullCmd
	}

	game := b.getGame(m.ChannelID)

	switch command {
	case "newgame":
		handleNewGame(s, m, game)
	case "join":
		handleJoin(s, m, game)
	case "start":
		handleStart(s, m, game)
	case "fold":
		handleFold(s, m, game)
	case "call":
		handleCall(s, m, game)
	case "raise":
		handleRaise(s, m, game, args)
	case "check":
		handleCheck(s, m, game)
	case "help":
		handleHelp(s, m)
	case "buyin":
		handleBuyIn(s, m, game, args)
	case "deal":
		handleDeal(s, m, game)
	case "count":
		handleCount(s, m, game)
	case "allin":
		handleAllIn(s, m, game)
	case "endgame":
		handleEndGame(s, m, game)
	case "options":
		handleOptions(s, m, game, args)
	}
}

func handleNewGame(s *discordgo.Session, m *discordgo.MessageCreate, game *Game) {
	if game.GetState() != NoGame {
		s.ChannelMessageSend(m.ChannelID, "A game is already in progress!")
		return
	}

	game.StartNewGame()
	s.ChannelMessageSend(m.ChannelID, "New game started! Type !join to join the game.")
}

func handleJoin(s *discordgo.Session, m *discordgo.MessageCreate, game *Game) {
	if game.GetState() != Waiting {
		s.ChannelMessageSend(m.ChannelID, "No game is waiting for players!")
		return
	}

	game.AddPlayer(m.Author)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s has joined the game!", m.Author.Username))
}

func handleStart(s *discordgo.Session, m *discordgo.MessageCreate, game *Game) {
	if game.GetState() != Waiting {
		s.ChannelMessageSend(m.ChannelID, "No game is waiting to start!")
		return
	}

	if len(game.GetPlayers()) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Need at least 2 players to start!")
		return
	}

	messages := game.DealHands()
	for _, msg := range messages {
		s.ChannelMessageSend(m.ChannelID, msg)
	}
}

func handleFold(s *discordgo.Session, m *discordgo.MessageCreate, game *Game) {
	if game.GetState() < HandsDealt || game.GetState() > RiverDealt {
		s.ChannelMessageSend(m.ChannelID, "No hand in progress!")
		return
	}

	currentPlayer := game.GetCurrentPlayer()
	if currentPlayer == nil || currentPlayer.User.ID != m.Author.ID {
		s.ChannelMessageSend(m.ChannelID, "It's not your turn!")
		return
	}

	messages := game.Fold()
	for _, msg := range messages {
		s.ChannelMessageSend(m.ChannelID, msg)
	}
}

func handleCall(s *discordgo.Session, m *discordgo.MessageCreate, game *Game) {
	if game.GetState() < HandsDealt || game.GetState() > RiverDealt {
		s.ChannelMessageSend(m.ChannelID, "No hand in progress!")
		return
	}

	currentPlayer := game.GetCurrentPlayer()
	if currentPlayer == nil || currentPlayer.User.ID != m.Author.ID {
		s.ChannelMessageSend(m.ChannelID, "It's not your turn!")
		return
	}

	messages := game.Call()
	for _, msg := range messages {
		s.ChannelMessageSend(m.ChannelID, msg)
	}
}

func handleRaise(s *discordgo.Session, m *discordgo.MessageCreate, game *Game, args []string) {
	if game.GetState() < HandsDealt || game.GetState() > RiverDealt {
		s.ChannelMessageSend(m.ChannelID, "No hand in progress!")
		return
	}

	currentPlayer := game.GetCurrentPlayer()
	if currentPlayer == nil || currentPlayer.User.ID != m.Author.ID {
		s.ChannelMessageSend(m.ChannelID, "It's not your turn!")
		return
	}

	if len(args) != 1 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !raise <amount>")
		return
	}

	var amount int
	_, err := fmt.Sscanf(args[0], "%d", &amount)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Invalid amount!")
		return
	}

	messages := game.Raise(amount)
	for _, msg := range messages {
		s.ChannelMessageSend(m.ChannelID, msg)
	}
}

func handleCheck(s *discordgo.Session, m *discordgo.MessageCreate, game *Game) {
	if game.GetState() < HandsDealt || game.GetState() > RiverDealt {
		s.ChannelMessageSend(m.ChannelID, "No hand in progress!")
		return
	}

	currentPlayer := game.GetCurrentPlayer()
	if currentPlayer == nil || currentPlayer.User.ID != m.Author.ID {
		s.ChannelMessageSend(m.ChannelID, "It's not your turn!")
		return
	}

	messages := game.Check()
	for _, msg := range messages {
		s.ChannelMessageSend(m.ChannelID, msg)
	}
}

func handleBuyIn(s *discordgo.Session, m *discordgo.MessageCreate, game *Game, args []string) {
	if game.GetState() != Waiting {
		s.ChannelMessageSend(m.ChannelID, "No game is waiting for players!")
		return
	}

	if len(args) != 1 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !buyin <amount>")
		return
	}

	var amount int
	_, err := fmt.Sscanf(args[0], "%d", &amount)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Invalid amount!")
		return
	}

	messages := game.BuyIn(m.Author, amount)
	for _, msg := range messages {
		s.ChannelMessageSend(m.ChannelID, msg)
	}
}

func handleDeal(s *discordgo.Session, m *discordgo.MessageCreate, game *Game) {
	if game.GetState() != NoHands {
		s.ChannelMessageSend(m.ChannelID, "Cannot deal now!")
		return
	}

	messages := game.DealHands()
	for _, msg := range messages {
		s.ChannelMessageSend(m.ChannelID, msg)
	}
}

func handleCount(s *discordgo.Session, m *discordgo.MessageCreate, game *Game) {
	if game.GetState() == NoGame {
		s.ChannelMessageSend(m.ChannelID, "No game in progress!")
		return
	}

	players := game.GetPlayers()
	if len(players) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No players in the game!")
		return
	}

	status := "Player balances:"
	for _, p := range players {
		status += fmt.Sprintf("\n- %s: $%d", p.User.Username, p.Balance)
	}

	s.ChannelMessageSend(m.ChannelID, status)
}

func handleAllIn(s *discordgo.Session, m *discordgo.MessageCreate, game *Game) {
	if game.GetState() < HandsDealt || game.GetState() > RiverDealt {
		s.ChannelMessageSend(m.ChannelID, "No hand in progress!")
		return
	}

	currentPlayer := game.GetCurrentPlayer()
	if currentPlayer == nil || currentPlayer.User.ID != m.Author.ID {
		s.ChannelMessageSend(m.ChannelID, "It's not your turn!")
		return
	}

	messages := game.AllIn()
	for _, msg := range messages {
		s.ChannelMessageSend(m.ChannelID, msg)
	}
}

func handleEndGame(s *discordgo.Session, m *discordgo.MessageCreate, game *Game) {
	if game.GetState() == NoGame {
		s.ChannelMessageSend(m.ChannelID, "No game in progress!")
		return
	}

	messages := game.EndGame()
	for _, msg := range messages {
		s.ChannelMessageSend(m.ChannelID, msg)
	}
}

func handleOptions(s *discordgo.Session, m *discordgo.MessageCreate, game *Game, args []string) {
	if len(args) == 0 {
		// Show current options
		options := game.GetOptions()
		status := fmt.Sprintf("Current game options:\n"+
			"Small Blind: $%d\n"+
			"Big Blind: $%d\n"+
			"Min Buy-In: $%d\n"+
			"Max Buy-In: $%d\n"+
			"Blind Raise Delay: %d minutes (0 = off)",
			options.SmallBlind, options.BigBlind, options.MinBuyIn, options.MaxBuyIn, options.RaiseDelay)
		s.ChannelMessageSend(m.ChannelID, status)
		return
	}

	if game.GetState() != Waiting && game.GetState() != HandsDealt {
		s.ChannelMessageSend(m.ChannelID, "Can only set options between hands!")
		return
	}

	if len(args) != 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !options [sb|bb|min|max|delay] <amount>")
		return
	}

	option := strings.ToLower(args[0])
	amount, err := strconv.Atoi(args[1])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Invalid amount!")
		return
	}

	options := game.GetOptions()
	switch option {
	case "sb":
		if amount <= 0 {
			s.ChannelMessageSend(m.ChannelID, "Small blind must be greater than 0!")
			return
		}
		if amount >= options.BigBlind {
			s.ChannelMessageSend(m.ChannelID, "Small blind must be less than big blind!")
			return
		}
		options.SmallBlind = amount
	case "bb":
		if amount <= options.SmallBlind {
			s.ChannelMessageSend(m.ChannelID, "Big blind must be greater than small blind!")
			return
		}
		options.BigBlind = amount
	case "min":
		if amount <= 0 {
			s.ChannelMessageSend(m.ChannelID, "Min buy-in must be greater than 0!")
			return
		}
		if amount >= options.MaxBuyIn {
			s.ChannelMessageSend(m.ChannelID, "Min buy-in must be less than max buy-in!")
			return
		}
		options.MinBuyIn = amount
	case "max":
		if amount <= options.MinBuyIn {
			s.ChannelMessageSend(m.ChannelID, "Max buy-in must be greater than min buy-in!")
			return
		}
		options.MaxBuyIn = amount
	case "delay":
		if amount < 0 {
			s.ChannelMessageSend(m.ChannelID, "Delay must be 0 or greater!")
			return
		}
		options.RaiseDelay = amount
	default:
		s.ChannelMessageSend(m.ChannelID, "Invalid option! Use sb, bb, min, max, or delay")
		return
	}

	game.SetOptions(options)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s set to %d", option, amount))
}

func handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	help := `Available commands:
!newgame - Start a new game
!join - Join the current game
!buyin <amount> - Buy in with specified amount
!start - Start the game with current players
!deal - Deal the cards
!fold - Fold your hand
!call - Call the current bet
!raise <amount> - Raise the bet
!allin - Go all in
!check - Check if no bet is required
!count - Show player balances
!options [sb|bb|min|max|delay] <amount> - Show or set game options
!endgame - End the current game
!help - Show this help message`

	s.ChannelMessageSend(m.ChannelID, help)
}
