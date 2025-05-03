package main

import (
	"os"

	bot "go-poker-bot/Bot"
)

func main() {
	bot.BotToken = os.Getenv("DISCORD_BOT_TOKEN")
	bot.Run()
}
