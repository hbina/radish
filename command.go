package redis

// Command flags. Please check the command table defined in the redis.c file
// for more information about the meaning of every flag.
const (
	CMD_WRITE    uint64 = 1 << 0
	CMD_READONLY        = 1 << 1
)

const (
	BCMD_OK = iota
	BCMD_RETRY
)

type CommandHandler func(c *Client, cmd [][]byte)
type BlockingCommandHandler func(c *Client, cmd [][]byte) int

type Command struct {
	name    string
	handler CommandHandler
	flag    uint64
}

func NewCommand(name string, handler CommandHandler, flag uint64) *Command {
	return &Command{
		name:    name,
		handler: handler,
		flag:    flag,
	}
}

type BlockingCommand struct {
	name    string
	handler BlockingCommandHandler
	flag    uint64
}

// Gets registered commands name.
func (cmd *Command) Name() string {
	return cmd.name
}

func NewBlockingCommand(name string, handler BlockingCommandHandler, flag uint64) *BlockingCommand {
	return &BlockingCommand{
		name:    name,
		handler: handler,
		flag:    flag,
	}
}

type BlockedCommand struct {
	c    *Client
	args [][]byte
}

func (r *Redis) RegisterCommands(cmds []*Command) {
	for _, cmd := range cmds {
		r.commands[cmd.name] = cmd
	}
}

func (r *Redis) RegisterBlockingCommands(cmds []*BlockingCommand) {
	for _, cmd := range cmds {
		r.blockingCommands[cmd.name] = cmd
	}
}
