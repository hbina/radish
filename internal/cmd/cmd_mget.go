package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
)

// https://redis.io/commands/mget/
// MGET key [key ...]
func MgetCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()
	keys := make([]string, 0)

	for i := 1; i < len(args); i++ {
		keys = append(keys, string(args[i]))
	}

	c.Conn().WriteArray(len(keys))
	for _, key := range keys {
		maybeItem, _ := db.GetOrExpire(key, true)

		if maybeItem == nil || maybeItem.Type() != types.ValueTypeString {
			c.Conn().WriteNull()
		} else {
			item := maybeItem.(*types.String)
			c.Conn().WriteBulkString(item.AsString())
		}
	}
}
