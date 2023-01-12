package cmd

import (
	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/type/
// TYPE key
func TypeCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 2 {
		c.WriteError(util.WrongNumOfArgsErr)
		return
	}

	key := string(args[1])
	db := c.Db()

	maybeItem, _ := db.Get(key)

	if maybeItem == nil {
		c.WriteSimpleString("none")
	} else {
		c.WriteSimpleString(maybeItem.TypeFancy())
	}
}
