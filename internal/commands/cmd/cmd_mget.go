package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/mget/
// MGET key [key ...]
func MgetCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()
	keys := make([]string, 0)

	for i := 1; i < len(args); i++ {
		keys = append(keys, string(args[i]))
	}

	c.Conn().WriteArray(len(keys))
	for _, key := range keys {
		maybeItem, _ := db.Get(key)

		if maybeItem == nil || maybeItem.Type() != types.ValueTypeString {
			if c.R3 {
				c.Conn().WriteNull()
			} else {
				c.Conn().WriteNullBulk()
			}
		} else {
			item := maybeItem.(*types.String)
			c.Conn().WriteBulkString(item.AsString())
		}
	}
}
