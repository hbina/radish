package redis

import (
	"fmt"
)

// https://redis.io/commands/mget/
// MGET key [key ...]
func StrlenCommand(c *Client, args [][]byte) {
	if len(args) != 2 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()
	key := string(args[1])

	maybeItem, _ := db.GetOrExpire(key, true)

	if maybeItem == nil {
		c.Conn().WriteInt(0)
	} else if maybeItem.Type() != ValueTypeString {
		c.Conn().WriteError(WrongTypeErr)
	} else {
		item := maybeItem.(*String)
		c.Conn().WriteInt(item.Len())
	}
}
