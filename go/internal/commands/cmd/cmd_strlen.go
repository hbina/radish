package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/strlen/
// STRLEN key
func StrlenCommand(c *pkg.Client, args [][]byte) {
	if len(args) != 2 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()
	key := string(args[1])

	maybeItem, _ := db.Get(key)

	if maybeItem == nil {
		c.Conn().WriteInt(0)
		return
	} else if maybeItem.Type() != types.ValueTypeString {
		c.Conn().WriteError(util.WrongTypeErr)
		return
	} else {
		item := maybeItem.(*types.String)
		c.Conn().WriteInt(item.Len())
		return
	}
}
