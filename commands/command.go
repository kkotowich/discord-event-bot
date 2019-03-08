package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type action func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error

// Command text map to function or other sub commands
type Command struct {
	//subcommands
	subcommands map[string]*Command

	Name        string
	Description string

	commandAction action
}

// AddSubCommand set passed command as to be a subcommand
func (c *Command) AddSubCommand(subcommand *Command) {
	if c.subcommands == nil {
		c.subcommands = make(map[string]*Command)
	}
	c.subcommands[subcommand.Name] = subcommand
}

// SetCommandAction set the action function for this command
func (c *Command) SetCommandAction(a action) {
	c.commandAction = a
}

func (c *Command) help(s *discordgo.Session, m *discordgo.MessageCreate, err error) {
	var sb strings.Builder

	// command's name and description
	sb.WriteString(c.Name)
	sb.WriteRune('\n')
	sb.WriteString(c.Description)
	sb.WriteRune('\n')

	// subcommand list
	for k, v := range c.subcommands {
		sb.WriteString(k)
		sb.WriteString(":\t")
		sb.WriteString(v.Description)
		sb.WriteRune('\n')
	}

	// display error
	if err != nil {
		sb.WriteString("error:\n```")
		sb.WriteString(err.Error())
		sb.WriteString("```")
	}
	s.ChannelMessageSend(m.ChannelID, sb.String())
}

// RunCommand runs action for this command, handles passing to sub command if they exist instead.
// When an error occurs or no subcommand exist, display commands help text
func (c *Command) RunCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	// check if subcommand exists first
	if len(c.subcommands) > 0 {
		if len(args) > 0 {
			subcommand, exists := c.subcommands[args[0]]
			if exists {
				var subargs []string
				if len(args) > 1 {
					subargs = args[1:]
				}
				subcommand.RunCommand(s, m, subargs)
			} else {
				c.help(s, m, nil)
			}
		} else {
			c.help(s, m, nil)
		}
	} else {
		err := c.commandAction(s, m, args)
		if err != nil {
			c.help(s, m, err)
		}
	}
}
