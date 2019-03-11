package commands

import (
	"discord-event-bot/types"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type dateDiffArgs struct {
	StartDate types.Date
	EndDate   types.Date
}

func (s dateDiffArgs) NewStruct() argsInterface {
	return &dateDiffArgs{}
}

func setDateDifCommand() {
	dateDiff := Command{Name: "datediff", Description: "get the difference in dates between two dates"}

	dateDiff.SetCommandAction(dateDiffAction, &dateDiffArgs{})

	Commands[dateDiff.Name] = &dateDiff
}

func dateDiffAction(s *discordgo.Session, m *discordgo.MessageCreate, argStruct argsInterface) error {
	d := argStruct.(dateDiffArgs)

	diff := d.EndDate.Date.Sub(d.StartDate.Date)

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%.2f", diff.Hours()/24))

	return nil
}
