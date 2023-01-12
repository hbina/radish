package cmd

import "github.com/hbina/radish/internal/pkg"

// https://redis.io/commands/exec/
// EXEC
func ExecCommand(c *pkg.Client, args [][]byte) {
	// Currently no-op because we are not multi-threaded to begin with
	c.WriteSimpleString("OK")
}
