package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Ready(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Printf("event.Version: %v\n", event.Version)
}
