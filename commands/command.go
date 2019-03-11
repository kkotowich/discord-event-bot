package commands

import (
	"discord-event-bot/services"
	"discord-event-bot/types"
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type argsInterface interface {
	NewStruct() argsInterface
}

type action func(s *discordgo.Session, m *discordgo.MessageCreate, argStruct argsInterface) error

// Command text map to function or other sub commands
type Command struct {
	// for help text to get full command
	Parent *Command

	//subcommands
	subcommands map[string]*Command

	Name        string
	Description string

	commandAction action

	// helper to validate arguments against
	argStruct        argsInterface
	argRequiredCount int
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
func (c *Command) SetCommandAction(a action, argStruct argsInterface) {
	c.commandAction = a
	if argStruct != nil {
		c.argStruct = argStruct
		c.argRequiredCount = reflect.ValueOf(c.argStruct).Elem().NumField()
	}
}

// SetRequiredCount set the required arg count, the remaining is considered optional
func (c *Command) SetRequiredCount(count int) {
	c.argRequiredCount = count
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
	if c.argStruct == nil {
		return
	}
	s := reflect.ValueOf(c.argStruct).Elem()

	// loop through each field and display its name and type
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		sb.WriteString(" <")
		sb.WriteString(s.Type().Field(i).Name)
		switch f.Type().String() {

		case "int":
			sb.WriteString(" Integer")
		case "string":
			sb.WriteString(" String")
		case "[]string":
			sb.WriteString("...list: Strings")
		case "[]int":
			sb.WriteString("...list: Integers")
		case "time.Time":
			sb.WriteString(" DateTime")
		case "types.Date":
			sb.WriteString(" Date")
		case "types.MilitaryTime":
			sb.WriteString(" MilitaryTime")
		default:
			sb.WriteString(" **UNKNOWN TYPE!**")
		}

		if i+1 > c.argRequiredCount {
			sb.WriteString(" (optional)")
		}
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

// convert args to struct
func (c *Command) convertArgs(args []string) (argsInterface, error) {
	var lastIsArray bool
	var sb strings.Builder

	// no argStruct
	if c.argStruct == nil {
		return nil, nil
	}

	s := reflect.ValueOf(c.argStruct.NewStruct()).Elem()
	totalFields := s.NumField()
	totalArgs := len(args)

	// validate correct amount of arguments is passed
	if totalArgs > totalFields {
		// if more args exist this is ok if the last variable in the struct is a slice
		if s.Field(totalFields-1).Kind() == reflect.Slice {
			lastIsArray = true
		} else {
			return nil, errors.New("Too many arguments were given")
		}
	} else if totalArgs < c.argRequiredCount {
		return nil, errors.New("Not enough arguments were given")
	}

	// only process upto totalArgs if less than total fields
	if totalArgs < totalFields {
		totalFields = totalArgs
	}

	// loop through each field and set the args
	for i := 0; i < totalFields; i++ {
		f := s.Field(i)

		// convert from string to one of the supported types
		switch f.Type().String() {

		case "int":
			intValue, err := strconv.Atoi(args[i])
			if err != nil {
				sb.WriteString(s.Type().Field(i).Name)
				sb.WriteString(" - Unable to convert \"")
				sb.WriteString(args[i])
				sb.WriteString("\" to a number")
				return nil, errors.New(sb.String())
			}
			f.SetInt(int64(intValue))

		case "string":
			f.SetString(args[i])
		case "[]string":
			var sliceArgs []string

			if lastIsArray && i == s.NumField()-1 {
				sliceArgs = args[i:]
			} else {
				sliceArgs = strings.Split(args[i], ",")
			}
			services.TrimSlice(&sliceArgs)

			f.Set(reflect.ValueOf(sliceArgs))

		case "[]int":
			var sliceArgs []string

			if lastIsArray && i == s.NumField()-1 {
				sliceArgs = args[i:]
			} else {
				sliceArgs = strings.Split(args[i], ",")
			}
			services.TrimSlice(&sliceArgs)

			intArgs, err := services.ConvertIntSlice(sliceArgs)
			if err != nil {
				sb.WriteString(s.Type().Field(i).Name)
				sb.WriteString(" - ")
				sb.WriteString(err.Error())
				return nil, errors.New(sb.String())
			}
			f.Set(reflect.ValueOf(intArgs))
		case "time.Time":
			layout := "2006-01-02 15:04"

			t, err := time.Parse(layout, args[i])
			if err != nil {
				sb.WriteString(s.Type().Field(i).Name)
				sb.WriteString(" - Unable to convert datetime: ")
				sb.WriteString(args[i])
				return nil, errors.New(sb.String())
			}
			f.Set(reflect.ValueOf(t))
		case "types.Date":
			layout := "2006-01-02"

			t, err := time.Parse(layout, args[i])
			if err != nil {
				sb.WriteString(s.Type().Field(i).Name)
				sb.WriteString(" - Unable to convert Date: ")
				sb.WriteString(args[i])
				return nil, errors.New(sb.String())
			}
			date := types.Date{Date: t}
			f.Set(reflect.ValueOf(date))
		case "types.MilitaryTime":
			layout := "15:04"

			t, err := time.Parse(layout, args[i])
			if err != nil {
				sb.WriteString(s.Type().Field(i).Name)
				sb.WriteString(" - Unable to convert Time: ")
				sb.WriteString(args[i])
				return nil, errors.New(sb.String())
			}
			militaryTime := types.MilitaryTime{Time: t}
			f.Set(reflect.ValueOf(militaryTime))
		default:
			sb.WriteString("Implementation Error - Unsupported Argument Type: ")
			sb.WriteString(f.Type().String())
			return nil, errors.New(sb.String())
		}
	}

	return s.Interface().(argsInterface), nil
}

// validate args then run
func (c *Command) validateAndRun(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	i, err := c.convertArgs(args)
	if err != nil {
		return err
	}
	return c.commandAction(s, m, i)
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
		err := c.validateAndRun(s, m, args)
		if err != nil {
			c.Help(s, m, err)
		}
	}
}
