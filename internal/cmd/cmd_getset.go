package cmd

import (
	"fmt"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
)

// https://redis.io/commands/getset/
// GETSET key value
// Note that this command is due for deprecation
func GetsetCommand(c *pkg.Client, args [][]byte) {
	if len(args) != 3 {
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()
	key := string(args[1])
	value := string(args[2])

	maybeItem, _ := db.GetOrExpire(key, true)

	if maybeItem != nil && maybeItem.Type() != types.ValueTypeString {
		c.Conn().WriteError(pkg.WrongTypeErr)
		return
	}

	db.Set(key, types.NewString(value), time.Time{})

	if maybeItem == nil {
		c.Conn().WriteNull()
	} else {
		// We already asserted that maybeItem is not nil and that it is a string
		c.Conn().WriteBulkString(maybeItem.(*types.String).Inner)
	}
}
