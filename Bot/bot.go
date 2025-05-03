package Bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func checkNilErr(e error) {
	if e != nil {
		log.Fatal("Error message")
	}
}

func Run(token string) {
	discord, err := discordgo.New("Bot " + token)
	checkNilErr(err)

	discord.AddHandler(newMessage)

	discord.Open()
	defer discord.Close()

	fmt.Println("Bot running....")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.Bot {
		return
	}

	switch {
	case strings.Contains(message.Content, "!help"):
		discord.ChannelMessageSend(message.ChannelID, "Hello World😃")
	case strings.Contains(message.Content, "!bye"):
		discord.ChannelMessageSend(message.ChannelID, "Good Bye👋")
		// add more cases if required
	}
}
