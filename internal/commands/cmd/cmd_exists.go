package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/exists/
func ExistsCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()
	count := 0

	for i := 1; i < len(args); i++ {
		key := string(args[i])
		value, _ := db.Get(key)
		if value != nil {
			count++
		}
	}

	c.Conn().WriteInt(count)
}
