package redis

import "github.com/redis-go/redcon"

// Command flags. Please check the command table defined in the redis.c file
// for more information about the meaning of every flag.
const (
	CMD_WRITE    uint64 = 1 << 0 /* "w" flag */
	CMD_READONLY        = 1 << 1 /* "r" flag */
	CMD_DENYOOM         = 1 << 2 /* "m" flag */
	CMD_MODULE          = 1 << 3 /* Command exported by module. */
	CMD_ADMIN           = 1 << 4 /* "a" flag */
	CMD_PUBSUB          = 1 << 5 /* "p" flag */
	CMD_NOSCRIPT        = 1 << 6 /* "s" flag */
	CMD_BLOCKING        = 1 << 8
	CMD_LOADING         = 1 << 9
	CMD_STALE           = 1 << 10
	CMD_FAST            = 1 << 14
	// Add more commands whenever necessary
)

// A command can be registered.
type Command struct {
	// The command name.
	name string

	// Handler
	handler CommandHandler

	// Command flag
	flag uint64 // Use map as a set data structure
}

func NewCommand(name string, handler CommandHandler, flags ...uint64) *Command {
	var flag uint64 = 0
	for _, f := range flags {
		flag = flag | f
	}

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
type CommandHandler func(c *Client, cmd redcon.Command)

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
	r.Mu().Lock()
	defer r.Mu().Unlock()
	for _, cmd := range cmds {
		r.registerCommand(cmd)
	}
}

// RegisterCommand adds a command to the redis instance.
// If cmd already exists the handler is overridden.
func (r *Redis) RegisterCommand(cmd *Command) {
	r.Mu().Lock()
	defer r.Mu().Unlock()
	r.registerCommand(cmd)
}
func (r *Redis) registerCommand(cmd *Command) {
	r.getCommands()[cmd.Name()] = cmd
}

// UnregisterCommand removes a command.
func (r *Redis) UnregisterCommand(name string) {
	r.Mu().Lock()
	defer r.Mu().Unlock()
	delete(r.commands, name)
}

// Command returns the registered command or nil if not exists.
func (r *Redis) Command(name string) *Command {
	r.Mu().RLock()
	defer r.Mu().RUnlock()
	return r.command(name)
}

func (r *Redis) command(name string) *Command {
	return r.commands[name]
}

// Commands returns the commands map.
func (r *Redis) Commands() Commands {
	r.Mu().RLock()
	defer r.Mu().RUnlock()
	return r.getCommands()
}

func (r *Redis) getCommands() Commands {
	return r.commands
}

// CommandExists checks if one or more commands are registered.
func (r *Redis) CommandExists(cmds ...string) bool {
	regCmds := r.Commands()

	for _, cmd := range cmds {
		if _, ex := regCmds[cmd]; !ex {
			return false
		}
	}
	return true
}

// FlushCommands removes all commands.
func (r *Redis) FlushCommands() {
	r.Mu().Lock()
	defer r.Mu().Unlock()
	r.commands = make(Commands)
}

// CommandHandlerFn returns the CommandHandler of cmd.
func (r *Redis) CommandHandlerFn(name string) CommandHandler {
	r.Mu().RLock()
	defer r.Mu().RUnlock()
	return r.command(name).handler
}

// UnknownCommandFn returns the UnknownCommand function.
func (r *Redis) UnknownCommandFn() UnknownCommand {
	r.Mu().RLock()
	defer r.Mu().RUnlock()
	return r.unknownCommand
}
