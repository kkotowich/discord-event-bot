package commands

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Commands mapping of top level commands
var Commands map[string]*Command

// SetCommands set all command mappings
func SetCommands() {
	Commands = make(map[string]*Command)

	// example ping command
	ping := Command{Name: "ping", Description: "respond to ping"}

	ping.SetCommandAction(func(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) error {
		s.ChannelMessageSend(m.ChannelID, "pong!")
		return nil
	})
	Commands[ping.Name] = &ping

	// say command
	say := Command{Name: "say", Description: "bot will repeat what you said"}
	Commands[say.Name] = &say

	// say subcommands
	sayUpper := Command{Name: "upper", Description: "say it in upper case"}
	sayUpper.SetCommandAction(func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
		if len(args) == 0 {
			return errors.New("Pass something to say")
		}
		s.ChannelMessageSend(m.ChannelID, strings.ToUpper(strings.Join(args, " ")))
		return nil
	})

	sayLower := Command{Name: "lower", Description: "say it in lower case"}
	sayLower.SetCommandAction(func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
		if len(args) == 0 {
			return errors.New("Pass something to say")
		}
		s.ChannelMessageSend(m.ChannelID, strings.ToLower(strings.Join(args, " ")))
		return nil
	})

	// add subcommands
	say.AddSubCommand(&sayUpper)
	say.AddSubCommand(&sayLower)

}
