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
		c.WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()
	keys := make([]string, 0)

	for i := 1; i < len(args); i++ {
		keys = append(keys, string(args[i]))
	}

	c.WriteArray(len(keys))
	for _, key := range keys {
		maybeItem, _ := db.Get(key)

		if maybeItem == nil || maybeItem.Type() != types.ValueTypeString {
			c.WriteNull()
		} else {
			item := maybeItem.(*types.String)
			c.WriteBulkString(item.AsString())
		}
	}
}
