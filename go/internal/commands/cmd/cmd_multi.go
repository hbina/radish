package cmd

import "github.com/hbina/radish/internal/pkg"

// https://redis.io/commands/multi/
// MULTI
func MultiCommand(c *pkg.Client, args [][]byte) {
	// Currently no-op because we are not multi-threaded to begin with
	c.Conn().WriteString("OK")
}
