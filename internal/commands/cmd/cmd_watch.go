package cmd

import (
	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/watch/
// WATCH key [key ...]
func WatchCommand(c *pkg.Client, args [][]byte) {
	if len(args) == 0 {
		c.WriteError(util.ZeroArgumentErr)
		return
	}

	// Currently no-op because we are not multi-threaded to begin with
	c.WriteSimpleString("OK")
}
