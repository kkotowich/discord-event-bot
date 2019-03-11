package commands

import (
	"errors"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

type counterArgs struct {
	Start     int
	End       int
	Increment int
}

func (s counterArgs) NewStruct() argsInterface {
	return &counterArgs{}
}

func setCounter() {
	counter := Command{Name: "counter", Description: "sum numbers in a range"}

	counter.SetCommandAction(counterAction, &counterArgs{})

	Commands[counter.Name] = &counter
}

func counterAction(s *discordgo.Session, m *discordgo.MessageCreate, argStruct argsInterface) error {
	c := argStruct.(counterArgs)
	// convert args to int
	if c.End < c.Start {
		return errors.New("end must be greater than start")
	}

	if c.Increment <= 0 {
		return errors.New("increment needs to be greater than 0")
	}

	var total int
	for i := c.Start; i <= c.End; i += c.Increment {
		total += i
	}
	s.ChannelMessageSend(m.ChannelID, strconv.Itoa(total))
	return nil
}
