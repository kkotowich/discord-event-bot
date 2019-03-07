package handlers

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func parseCommand(s string) (string, []string) {
	var fields []string
	var sb strings.Builder
	var quote rune

	// parse string, spliting on spaces tabs and newlines
	// get everything in quotes
	for _, c := range s {
		if quote == 0 && (c == '\'' || c == '"' || c == '`') && sb.Len() == 0 {
			quote = c
		} else if quote == c || (quote == 0 && (c == ' ' || c == '\n' || c == '\t')) {
			quote = 0

			field := sb.String()
			if len(field) > 0 {
				fields = append(fields, field)
				sb.Reset()
			}
		} else {
			sb.WriteRune(c)
		}
	}

	// write out last string
	field := sb.String()
	if len(field) > 0 {
		fields = append(fields, field)
	}

	return fields[0], fields[1:]
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Ignore if message does not start with configured prefix
	prefix := os.Getenv("CommandPrefix")
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	// Split command and args
	text := strings.TrimSpace(strings.TrimPrefix(m.Content, prefix))
	if len(text) == 0 {
		return
	}

	command, args := parseCommand(text)

	// TODO: map\register command to function and call here with supplied args
	// i.e. commands[command](args)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("command: %s\nargs: %s", command, strings.Join(args, ", ")))
}
