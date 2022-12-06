package redis

import "github.com/tidwall/redcon"

// Command flags. Please check the command table defined in the redis.c file
// for more information about the meaning of every flag.
const (
	CMD_WRITE    uint64 = 1 << 0
	CMD_READONLY        = 1 << 1
)

// A command can be registered.
type Command struct {
	// The command name.
	name string

	// Handler
	handler CommandHandler

	// Command flag
	flag uint64
}

func NewCommand(name string, handler CommandHandler, flag uint64) *Command {

	return &Command{
		name:    name,
		handler: handler,
		flag:    flag,
	}
}

// Command flag type.
type CmdFlag uint

// Commands map
type Commands map[string]*Command

// The CommandHandler is triggered when the received
// command equals a registered command.
//
// However the CommandHandler is executed by the Handler,
// so if you implement an own Handler make sure the CommandHandler is called.
type CommandHandler func(c *Client, cmd [][]byte)

// Is called when a request is received,
// after Accept and if the command is not registered.
//
// However UnknownCommand is executed by the Handler,
// so if you implement an own Handler make sure to include UnknownCommand.
type UnknownCommand func(c *Client, cmd redcon.Command)

// Gets registered commands name.
func (cmd *Command) Name() string {
	return cmd.name
}

// RegisterCommands adds commands to the redis instance.
// If a cmd already exists the handler is overridden.
func (r *Redis) RegisterCommands(cmds []*Command) {

	for _, cmd := range cmds {
		r.commands[cmd.Name()] = cmd
	}
}

// Command returns the registered command or nil if not exists.
func (r *Redis) Command(name string) *Command {

	return r.commands[name]
}

// Commands returns the commands map.
func (r *Redis) Commands() Commands {

	return r.commands
}

// UnknownCommandFn returns the UnknownCommand function.
func (r *Redis) UnknownCommandFn() UnknownCommand {

	return r.unknownCommand
}
