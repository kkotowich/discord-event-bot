package commands

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type sayArgs struct {
	Words []string
}

func (s sayArgs) NewStruct() argsInterface {
	return &sayArgs{}
}

func setSayCommand() {
	// say command
	say := Command{Name: "say", Description: "bot will repeat what you said"}
	Commands[say.Name] = &say

	// say subcommands
	sayUpper := Command{Name: "upper", Description: "say it in upper case"}
	sayUpper.SetCommandAction(sayUpperAction, &sayArgs{})

	sayLower := Command{Name: "lower", Description: "say it in lower case"}
	sayLower.SetCommandAction(sayLowerAction, &sayArgs{})

	// add subcommands
	say.AddSubCommand(&sayUpper)
	say.AddSubCommand(&sayLower)
}

func sayUpperAction(s *discordgo.Session, m *discordgo.MessageCreate, argStruct argsInterface) error {
	sa := argStruct.(sayArgs)
	if len(sa.Words) == 0 {
		return errors.New("Pass something to say")
	}
	s.ChannelMessageSend(m.ChannelID, strings.ToUpper(strings.Join(sa.Words, " ")))
	return nil
}

func sayLowerAction(s *discordgo.Session, m *discordgo.MessageCreate, argStruct argsInterface) error {
	sa := argStruct.(sayArgs)
	if len(sa.Words) == 0 {
		return errors.New("Pass something to say")
	}
	s.ChannelMessageSend(m.ChannelID, strings.ToLower(strings.Join(sa.Words, " ")))
	return nil
}
