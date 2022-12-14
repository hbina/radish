package cmd

import "github.com/hbina/radish/internal/pkg"

// https://redis.io/commands/object/
func ObjectCommand(c *pkg.Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError(pkg.ZeroArgumentErr)
		return
	}

	c.Conn().WriteNull()
}
