# Go Poker Bot

A Discord poker bot written in Go.

## Installation

1. Install Go: [https://golang.org/doc/install](https://golang.org/doc/install)

2. Install dependencies:
```bash
go mod tidy
```

## Environment Setup

Create a `.env` file in the project root with your Discord bot token:
```
DISCORD_BOT_TOKEN=your_bot_token_here
```

To get your bot token:
1. Go to the [Discord Developer Portal](https://discord.com/developers/applications)
2. Create a new application
3. Go to the "Bot" tab
4. Click "Add Bot"
5. Copy the token and paste it in your `.env` file

The bot will automatically load the token from the `.env` file or default to your enviornment variables. 

Now, go to [this page](https://finitereality.github.io/permissions-calculator/?v=0), select all the Non-Administrative permissions, enter the client id from the bot's application page. Open the generated invite link and then select one of the servers you own to add it that server.

3. Run the bot
```bash
go run main.go
```

Finally, you can message `!newgame` to start playing. 