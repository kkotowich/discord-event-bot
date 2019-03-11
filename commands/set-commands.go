package commands

// Commands mapping of top level commands
var Commands map[string]*Command

// SetCommands set all command mappings
func SetCommands() {
	Commands = make(map[string]*Command)

	setHelpCommand()
	setPingCommand()
	setSayCommand()
	setCounter()
	setDateDifCommand()
}
