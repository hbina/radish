package cmd

import (
	"github.com/hbina/radish/internal/pkg"
)

// https://redis.io/commands/object/
func ObjectCommand(c *pkg.Client, args [][]byte) {
	if c.R3 {
		c.Conn().WriteNull()
	} else {
		c.Conn().WriteNullBulk()
	}
}
