package redis

import (
	"strings"
)

// https://redis.io/commands/flushall/
func FlushAllCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
	} else if len(args) == 1 || (len(args) == 2 && strings.ToLower(string(args[1])) == "sync") {
		syncCleanup(c)
		c.Conn().WriteString("OK")
	} else if len(args) == 2 && strings.ToLower(string(args[1])) == "async" {
		c.Conn().WriteError("FLUSHALL ASYNC is not implemented yet")
	} else {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
	}
}

func syncCleanup(c *Client) {
	c.Redis().SyncFlushAll()
}
