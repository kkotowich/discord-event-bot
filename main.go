package main

import (
	"discord-notify-bot/commands"
	"discord-notify-bot/handlers"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	SetEnvVars()
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + os.Getenv("BotAPIToken"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Set up commands
	commands.SetCommands()

	// Register handlers here
	dg.AddHandler(handlers.Ready)
	dg.AddHandler(handlers.MessageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
