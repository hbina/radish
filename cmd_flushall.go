package redis

import (
	"strings"

	"github.com/tidwall/redcon"
)

// https://redis.io/commands/flushall/
func FlushAllCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
	} else if len(cmd.Args) == 1 || (len(cmd.Args) == 2 && strings.ToLower(string(cmd.Args[1])) == "sync") {
		syncCleanup(c)
		c.Conn().WriteString("OK")
	} else if len(cmd.Args) == 2 && strings.ToLower(string(cmd.Args[1])) == "async" {
		c.Conn().WriteError("FLUSHALL ASYNC is not implemented yet")
	} else {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
	}
}

func syncCleanup(c *Client) {
	c.Redis().SyncFlushAll()
}
