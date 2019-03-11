package commands

import "github.com/bwmarrin/discordgo"

func setPingCommand() {
	// example ping command
	ping := Command{Name: "ping", Description: "respond to ping"}

	ping.SetCommandAction(func(s *discordgo.Session, m *discordgo.MessageCreate, _ argsInterface) error {
		s.ChannelMessageSend(m.ChannelID, "pong!")
		return nil
	}, nil)
	Commands[ping.Name] = &ping
}
