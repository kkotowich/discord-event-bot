package commands

import (
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type action func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error

// Command text map to function or other sub commands
type Command struct {
	// for help text to get full command
	Parent *Command

	//subcommands
	subcommands map[string]*Command

	Name        string
	Description string

	commandAction action

	// list of args to help identify what to enter
	argumentNames []string
}

// AddSubCommand set passed command as to be a subcommand
func (c *Command) AddSubCommand(subcommand *Command) {
	if c.subcommands == nil {
		c.subcommands = make(map[string]*Command)
	}
	subcommand.Parent = c
	c.subcommands[subcommand.Name] = subcommand
}

// SetCommandAction set the action function for this command
func (c *Command) SetCommandAction(a action, args []string) {
	c.commandAction = a
	c.argumentNames = args
}

// GetSubCommand return the subcommand in provide path
// nothing is return if not on path
func (c *Command) GetSubCommand(commands *[]string, i int) *Command {
	if c.subcommands == nil {
		return nil
	}

	subcommand, exists := c.subcommands[(*commands)[i]]
	if exists {
		if len(*commands) == i+1 {
			return subcommand
		}
		return subcommand.GetSubCommand(commands, i+1)
	}
	return nil
}

func (c *Command) commandText(sb *strings.Builder) {
	if c.Parent == nil {
		prefix := os.Getenv("CommandPrefix")
		sb.WriteString(prefix)
		sb.WriteString(c.Name)
	} else {
		c.Parent.commandText(sb)
		sb.WriteRune(' ')
		sb.WriteString(c.Name)
	}
}

func (c *Command) commandArgs(sb *strings.Builder) {
	for _, v := range c.argumentNames {
		sb.WriteString(" <")
		sb.WriteString(v)
		sb.WriteString(">")
	}
}

// Help provide descriptive help text for command
func (c *Command) Help(s *discordgo.Session, m *discordgo.MessageCreate, err error) {
	var sb strings.Builder

	// command's name and description
	sb.WriteString("`Command: ")
	c.commandText(&sb)
	c.commandArgs(&sb)
	sb.WriteString("`\n")
	sb.WriteString(c.Description)
	sb.WriteString("\n\n")
	// subcommand list
	if len(c.subcommands) > 0 {
		sb.WriteString("__**Subcommands:**__\n")
	}
	for k, v := range c.subcommands {
		sb.WriteString("__*")
		sb.WriteString(k)
		sb.WriteString(":*__ ")
		sb.WriteString(v.Description)
		sb.WriteRune('\n')
	}

	// display error
	if err != nil {
		sb.WriteString("__***error:***__\n```")
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
				c.Help(s, m, nil)
			}
		} else {
			c.Help(s, m, nil)
		}
	} else {
		err := c.commandAction(s, m, args)
		if err != nil {
			c.Help(s, m, err)
		}
	}
}
