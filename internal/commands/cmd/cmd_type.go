package cmd

import (
	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/type/
// TYPE key
func TypeCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 2 {
		c.Conn().WriteError(util.WrongNumOfArgsErr)
		return
	}

	key := string(args[1])
	db := c.Db()

	maybeItem, _ := db.GetOrExpire(key, true)

	if maybeItem == nil {
		c.Conn().WriteString("none")
	} else {
		c.Conn().WriteString(maybeItem.TypeFancy())
	}
}
