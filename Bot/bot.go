package Bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
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
		"pot": "allin",
	}

	// If command is a shorthand, replace it with the full version
	if fullCmd, exists := shorthand[command]; exists {
		command = fullCmd
	}

	game := b.getGame(m.ChannelID)

	// Lock the game for the duration of command processing
	game.mu.Lock()
	defer game.mu.Unlock()

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
	case "change":
		handleChangeGame(s, m, game, args)
	case "options":
		handleOptions(s, m, game, args)
	case "verbose":
		handleVerbose(s, m, game)
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

	if AddPlayer(s, m, game) {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s has joined the game!", m.Author.GlobalName))
		return
	}

	s.ChannelMessageSend(m.ChannelID, "You're already in the game!")
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

	SendMessages(s, m, game.DealHands())
	TellHands(s, m, game)
}

func handleFold(s *discordgo.Session, m *discordgo.MessageCreate, game *Game) {
	if game.GetState() < HandsDealt || game.GetState() > RiverDealt {
		s.ChannelMessageSend(m.ChannelID, "No hand in progress!")
		return
	}

	// Check if it's the current player's turn
	if !game.IsCurrentPlayer(m.Author) {
		s.ChannelMessageSend(m.ChannelID, "It's not your turn!")
		return
	}

	SendMessages(s, m, game.Fold())
}

func handleCall(s *discordgo.Session, m *discordgo.MessageCreate, game *Game) {
	if game.GetState() < HandsDealt || game.GetState() > RiverDealt {
		s.ChannelMessageSend(m.ChannelID, "No hand in progress!")
		return
	}

	// Check if it's the current player's turn
	if !game.IsCurrentPlayer(m.Author) {
		s.ChannelMessageSend(m.ChannelID, "It's not your turn!")
		return
	}

	SendMessages(s, m, game.Call())
}

func handleRaise(s *discordgo.Session, m *discordgo.MessageCreate, game *Game, args []string) {
	if game.GetState() < HandsDealt || game.GetState() > RiverDealt {
		s.ChannelMessageSend(m.ChannelID, "No hand in progress!")
		return
	}

	// Check if it's the current player's turn
	if !game.IsCurrentPlayer(m.Author) {
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

	SendMessages(s, m, game.Raise(amount))
}

func handleCheck(s *discordgo.Session, m *discordgo.MessageCreate, game *Game) {
	if game.GetState() < HandsDealt || game.GetState() > RiverDealt {
		s.ChannelMessageSend(m.ChannelID, "No hand in progress!")
		return
	}

	// Check if it's the current player's turn
	if !game.IsCurrentPlayer(m.Author) {
		s.ChannelMessageSend(m.ChannelID, "It's not your turn!")
		return
	}

	SendMessages(s, m, game.Check())
}

func handleBuyIn(s *discordgo.Session, m *discordgo.MessageCreate, game *Game, args []string) {
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

	newPlayer := AddPlayer(s, m, game)

	SendMessages(s, m, game.BuyIn(m.Author, amount, newPlayer))
}

func handleDeal(s *discordgo.Session, m *discordgo.MessageCreate, game *Game) {
	if game.GetState() != NoHands {
		s.ChannelMessageSend(m.ChannelID, "Cannot deal now!")
		return
	}

	SendMessages(s, m, game.DealHands())
	TellHands(s, m, game)
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
		status += fmt.Sprintf("\n- %s: $%d", p.Name, p.Balance)
	}

	s.ChannelMessageSend(m.ChannelID, status)
}

func handleAllIn(s *discordgo.Session, m *discordgo.MessageCreate, game *Game) {
	if game.GetState() < HandsDealt || game.GetState() > RiverDealt {
		s.ChannelMessageSend(m.ChannelID, "No hand in progress!")
		return
	}

	// Check if it's the current player's turn
	if !game.IsCurrentPlayer(m.Author) {
		s.ChannelMessageSend(m.ChannelID, "It's not your turn!")
		return
	}

	SendMessages(s, m, game.AllIn())
}

func handleEndGame(s *discordgo.Session, m *discordgo.MessageCreate, game *Game) {
	if game.GetState() == NoGame {
		s.ChannelMessageSend(m.ChannelID, "No game in progress!")
		return
	}

	SendMessages(s, m, game.EndGame())
}

func handleChangeGame(s *discordgo.Session, m *discordgo.MessageCreate, game *Game, args []string) {
	if !game.BetweenHands() {
		s.ChannelMessageSend(m.ChannelID, "Cannot change game type in the middle of a hand!")
		return
	}

	if len(args) != 1 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !change <holdem|plo>")
		return
	}

	gameType := args[0]
	message := game.ChangeGameType(gameType)
	s.ChannelMessageSend(m.ChannelID, message)
}

func handleOptions(s *discordgo.Session, m *discordgo.MessageCreate, game *Game, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, game.ListOptions())
		return
	}

	if !game.BetweenHands() {
		s.ChannelMessageSend(m.ChannelID, "Can only set options between hands!")
		return
	}

	if len(args) != 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !options [sb|bb|min|max|delay] <amount>")
		return
	}

	s.ChannelMessageSend(m.ChannelID, game.SetOption(args))
}

func handleVerbose(s *discordgo.Session, m *discordgo.MessageCreate, game *Game) {
	s.ChannelMessageSend(m.ChannelID, game.ToggleVerbose())
}

func TellHands(s *discordgo.Session, m *discordgo.MessageCreate, game *Game) {
	// for each player, send them a private message containing their dealt cards
	for _, player := range game.Players {
		channel, err := s.UserChannelCreate(player.User.ID)
		if err != nil {
			log.Fatal("Error fetching user:", err)
		}

		_, err = s.ChannelMessageSend(channel.ID, fmt.Sprintf("Your cards are: %s", player.PrintHand()))
		if err != nil {
			log.Fatal("Error sending DM message:", err)
			s.ChannelMessageSend(
				m.ChannelID,
				fmt.Sprintf("Failed to send %s a DM. Did you disable DM in your privacy settings?", player.Name),
			)
		}
	}
}

// Wrapper to set a player's nickname if it exists
func AddPlayer(s *discordgo.Session, m *discordgo.MessageCreate, game *Game) bool {
	if game.IsPlayer(m.Author) {
		return false
	}
	member, err := s.GuildMember(m.GuildID, m.Author.ID)
	name := m.Author.Username
	if err == nil && member.Nick != "" {
		name = member.Nick
	} else if m.Author.GlobalName != "" {
		name = m.Author.GlobalName
	}
	game.AddPlayer(m.Author, name)
	return true
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
!change <holdem|plo> - Change the game type
!help - Show this help message
!verbose - Toggle verbose output mode`

	s.ChannelMessageSend(m.ChannelID, help)
}

// SendMessages sends multiple messages to a channel with a delay between them
func SendMessages(s *discordgo.Session, m *discordgo.MessageCreate, messages []string) {
	for _, msg := range messages {
		s.ChannelMessageSend(m.ChannelID, msg)
	}
}
