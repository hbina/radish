package cmd

import (
	"strings"

	"github.com/hbina/radish/internal/pkg"
)

// https://redis.io/commands/flushall/
func FlushAllCommand(c *pkg.Client, args [][]byte) {
	if len(args) == 1 || (len(args) == 2 && strings.ToLower(string(args[1])) == "sync") {
		syncFlushAll(c)
		c.Conn().WriteString("OK")
	} else if len(args) == 2 && strings.ToLower(string(args[1])) == "async" {
		c.Conn().WriteError("FLUSHALL ASYNC is not implemented yet")
	} else {
		c.Conn().WriteString("OK")
	}
}

func syncFlushAll(c *pkg.Client) {
	c.Redis().SyncFlushAll()
}
