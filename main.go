package main

import (
	"os"

	bot "example.com/poker-bot/Bot"
)

func main() {
	bot.BotToken = os.Getenv("DISCORD_BOT_TOKEN")
	bot.Run()
}
