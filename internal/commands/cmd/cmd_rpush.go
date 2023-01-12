package cmd

import (
	"fmt"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/rpush/
// RPUSH key element [element ...]
func RPushCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	db := c.Db()
	maybeItem, _ := db.Get(key)

	if maybeItem == nil {
		maybeItem = types.NewList()
	} else if maybeItem.Type() != types.ValueTypeList {
		c.Conn().WriteError(util.WrongTypeErr)
		return
	}

	list := maybeItem.(*types.List)

	for i := 2; i < len(args); i++ {
		list.RPush(string(args[i]))
	}

	db.Set(key, list, time.Time{})

	c.Conn().WriteInt(list.Len())
}
