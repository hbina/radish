package cmd

import (
	"fmt"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
)

// https://redis.io/commands/rpush/
// RPUSH key element [element ...]
func RPushCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	db := c.Db()
	maybeItem, _ := db.GetOrExpire(key, true)

	if maybeItem == nil {
		maybeItem = NewList()
	} else if maybeItem.Type() != types.ValueTypeList {
		c.Conn().WriteError(pkg.WrongTypeErr)
		return
	}

	list := maybeItem.(*List)

	for i := 2; i < len(args); i++ {
		list.RPush(string(args[i]))
	}

	db.Set(key, list, time.Time{})

	c.Conn().WriteInt(list.Len())
}
