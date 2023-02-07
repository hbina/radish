package cmd

import "github.com/hbina/radish/internal/pkg"

// https://redis.io/commands/dbsize/
func DbSizeCommand(c *pkg.Client, args [][]byte) {
	db := c.Db()
	c.Conn().WriteInt(db.Len())
}
