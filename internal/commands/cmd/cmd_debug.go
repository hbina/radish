package cmd

import "github.com/hbina/radish/internal/pkg"

// https://redis.io/commands/debug/
func DebugCommand(c *pkg.Client, args [][]byte) {
	c.WriteSimpleString("Not implemented")
}
