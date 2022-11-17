package redis

import (
	"fmt"
	"time"
)

// https://redis.io/commands/getset/
// GETSET key value
// Note that this command is due for deprecation
func GetsetCommand(c *Client, args [][]byte) {
	if len(args) != 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()
	key := string(args[1])
	value := string(args[2])

	maybeItem, _ := db.GetOrExpire(key, true)

	if maybeItem != nil && maybeItem.Type() != ValueTypeString {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	db.Set(key, NewString(value), time.Time{})

	if maybeItem == nil {
		c.Conn().WriteNull()
	} else {
		// We already asserted that maybeItem is not nil and that it is a string
		c.Conn().WriteBulkString(maybeItem.(*String).inner)
	}
}
