package cmd

import "github.com/hbina/radish/internal/pkg"

// https://redis.io/commands/del/
func DelCommand(c *pkg.Client, args [][]byte) {
	db := c.Db()
	keys := make([]string, 0, len(args)-1)

	for i := 1; i < len(args); i++ {
		k := string(args[i])
		keys = append(keys, k)
	}

	count := db.Delete(keys...)
	c.Conn().WriteInt(count)
}
