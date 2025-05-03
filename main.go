package main

import (
	"log"
	"os"

	"go-poker-bot/Bot"

	"github.com/joho/godotenv"
)

func main() {
	// Try to load .env file, but don't fail if it doesn't exist
	_ = godotenv.Load()

	// Get bot token from environment (system env takes precedence over .env)
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_BOT_TOKEN not set in environment or .env file")
	}

	// Run the bot
	Bot.Run(token)
}
