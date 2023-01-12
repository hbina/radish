package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/mget/
// MGET key [key ...]
func StrlenCommand(c *pkg.Client, args [][]byte) {
	if len(args) != 2 {
		c.WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()
	key := string(args[1])

	maybeItem, _ := db.Get(key)

	if maybeItem == nil {
		c.WriteInt(0)
	} else if maybeItem.Type() != types.ValueTypeString {
		c.WriteError(util.WrongTypeErr)
	} else {
		item := maybeItem.(*types.String)
		c.WriteInt(item.Len())
	}
}
