package redis

import (
	"fmt"
	"time"
)

// https://redis.io/commands/rpush/
// RPUSH key element [element ...]
func RPushCommand(c *Client, args [][]byte) {
	if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	db := c.Db()
	maybeItem, _ := db.GetOrExpire(key, true)

	if maybeItem == nil {
		maybeItem = NewList()
	} else if maybeItem.Type() != ValueTypeList {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	list := maybeItem.(*List)

	for i := 2; i < len(args); i++ {
		list.RPush(string(args[i]))
	}

	db.Set(key, list, time.Time{})

	c.Conn().WriteInt(list.Len())
}
