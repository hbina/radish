package pkg

import "time"

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
type BlockingCommandHandler func(c *Client, cmd [][]byte) *BlockedCommand

type Command struct {
	Name    string
	Handler CommandHandler
	Flag    uint64
}

func NewCommand(name string, handler CommandHandler, flag uint64) *Command {
	return &Command{
		Name:    name,
		Handler: handler,
		Flag:    flag,
	}
}

type BlockingCommand struct {
	Name    string
	Handler BlockingCommandHandler
	Flag    uint64
}

func NewBlockingCommand(name string, handler BlockingCommandHandler, flag uint64) *BlockingCommand {
	return &BlockingCommand{
		Name:    name,
		Handler: handler,
		Flag:    flag,
	}
}

type BlockedCommand struct {
	c    *Client
	args [][]byte
	ttl  time.Time
}

func NewBlockedCommand(c *Client, args [][]byte, ttl time.Time) *BlockedCommand {
	return &BlockedCommand{
		c:    c,
		args: args,
		ttl:  ttl,
	}
}

func (b *BlockedCommand) Client() *Client {
	return b.c
}

func (b *BlockedCommand) Args() [][]byte {
	return b.args
}

func (b *BlockedCommand) Ttl() time.Time {
	return b.ttl
}
