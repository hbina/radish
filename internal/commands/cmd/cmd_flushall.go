package cmd

import (
	"strings"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/flushall/
func FlushAllCommand(c *pkg.Client, args [][]byte) {
	if len(args) == 1 || (len(args) == 2 && strings.ToLower(string(args[1])) == "sync") {
		syncFlushAll(c)
		c.WriteSimpleString("OK")
	} else if len(args) == 2 && strings.ToLower(string(args[1])) == "async" {
		c.WriteError("FLUSHALL ASYNC is not implemented yet")
	} else {
		c.WriteError(util.ZeroArgumentErr)
	}
}

func syncFlushAll(c *pkg.Client) {
	c.Redis().SyncFlushAll()
}
