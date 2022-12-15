package cmd

import "github.com/hbina/radish/internal/pkg"

// https://redis.io/commands/watch/
// WATCH key [key ...]
func WatchCommand(c *pkg.Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError(pkg.ZeroArgumentErr)
		return
	}

	// Currently no-op because we are not multi-threaded to begin with
	c.Conn().WriteString("OK")
}
