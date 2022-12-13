package cmd

import "github.com/hbina/radish/internal/pkg"

// https://redis.io/commands/function/
func FunctionCommand(c *pkg.Client, args [][]byte) {
	c.Conn().WriteString("Not implemented")
}
