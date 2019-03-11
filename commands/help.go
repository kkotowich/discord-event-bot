package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type helpArgs struct {
	Commands []string
}

func (s helpArgs) NewStruct() argsInterface {
	return &helpArgs{}
}

func setHelpCommand() {
	help := Command{Name: "help", Description: "List Commands"}

	help.SetCommandAction(helpAction, &helpArgs{})
	help.SetRequiredCount(0)

	Commands[help.Name] = &help
}

func helpAction(s *discordgo.Session, m *discordgo.MessageCreate, argStruct argsInterface) error {
	h := argStruct.(helpArgs)
	if len(h.Commands) > 0 {
		// provide detail help text for provided command
		helpCommand(s, m, h.Commands)
	} else {
		//list all commands
		helpList(s, m)
	}

	return nil
}

func helpCommand(s *discordgo.Session, m *discordgo.MessageCreate, commands []string) {
	var sb strings.Builder
	var helpCommand *Command

	c, exists := Commands[commands[0]]
	if exists {
		if len(commands) > 1 {
			helpCommand = c.GetSubCommand(&commands, 1)

		} else {
			helpCommand = c
		}
	}
	if helpCommand != nil {
		helpCommand.Help(s, m, nil)
	} else {
		sb.WriteString("**Command __*")
		sb.WriteString(strings.Join(commands, " "))
		sb.WriteString("*__ not found**")
		s.ChannelMessageSend(m.ChannelID, sb.String())
	}
}

func helpList(s *discordgo.Session, m *discordgo.MessageCreate) {
	var sb strings.Builder

	// list commands
	for k, v := range Commands {
		sb.WriteString("__*")
		sb.WriteString(k)
		sb.WriteString(":*__ ")
		sb.WriteString(v.Description)
		sb.WriteRune('\n')
	}
	s.ChannelMessageSend(m.ChannelID, sb.String())
}
