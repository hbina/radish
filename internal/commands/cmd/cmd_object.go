package cmd

import (
	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/object/
func ObjectCommand(c *pkg.Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError(util.ZeroArgumentErr)
		return
	}

	c.Conn().WriteNull()
}
