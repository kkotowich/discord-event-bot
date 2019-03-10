package commands

import (
	"errors"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Commands mapping of top level commands
var Commands map[string]*Command

// SetCommands set all command mappings
func SetCommands() {
	Commands = make(map[string]*Command)

	setHelpCommand()
	setPingCommand()
	setSayCommand()
	setCounter()
}

func setHelpCommand() {
	help := Command{Name: "help", Description: "List Commands"}

	help.SetCommandAction(helpAction, []string{"commands...(optional)"})

	Commands[help.Name] = &help
}

func helpAction(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) > 0 {
		// provide detail help text for provided command
		helpCommand(s, m, args)
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

func setPingCommand() {
	// example ping command
	ping := Command{Name: "ping", Description: "respond to ping"}

	ping.SetCommandAction(func(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) error {
		s.ChannelMessageSend(m.ChannelID, "pong!")
		return nil
	}, nil)
	Commands[ping.Name] = &ping
}

func setSayCommand() {
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
	}, []string{"words..."})

	sayLower := Command{Name: "lower", Description: "say it in lower case"}
	sayLower.SetCommandAction(func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
		if len(args) == 0 {
			return errors.New("Pass something to say")
		}
		s.ChannelMessageSend(m.ChannelID, strings.ToLower(strings.Join(args, " ")))
		return nil
	}, []string{"words..."})

	// add subcommands
	say.AddSubCommand(&sayUpper)
	say.AddSubCommand(&sayLower)
}

func setCounter() {
	counter := Command{Name: "counter", Description: "sum numbers in a range"}

	counter.SetCommandAction(counterAction, []string{"start", "end", "increment"})

	Commands[counter.Name] = &counter
}

func counterAction(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	// convert args to int
	if len(args) < 3 {
		return errors.New("Not enough arguments")
	}

	start, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("start must be a whole number")
	}

	end, err := strconv.Atoi(args[1])
	if err != nil {
		return errors.New("end must be a whole number")
	}
	if end < start {
		return errors.New("end must be greater than start")
	}

	increment, err := strconv.Atoi(args[2])
	if err != nil {
		return errors.New("increment must be a whole number")
	}
	if increment <= 0 {
		return errors.New("increment needs to be greater than 0")
	}

	var total int
	for i := start; i <= end; i += increment {
		total += i
	}
	s.ChannelMessageSend(m.ChannelID, strconv.Itoa(total))
	return nil
}
