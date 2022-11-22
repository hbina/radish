package redis

import (
	"fmt"
	"strconv"
	"time"
)

// https://redis.io/commands/decr/
func DecrCommand(c *Client, args [][]byte) {
	if len(args) == 1 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()

	key := string(args[1])

	item, exists := db.storage[key]

	if !exists {
		db.Set(key, NewString("-1"), time.Time{})
		c.conn.WriteInt64(-1)
		return
	}

	value, ok := item.Value().(*string)

	if !ok {
		c.conn.WriteError(InvalidIntErr)
		return
	}

	intValue, err := strconv.ParseInt(*value, 10, 64)

	if err != nil {
		c.conn.WriteError(WrongTypeErr)
		return
	}

	intValue--

	db.Set(key, NewString(fmt.Sprint(intValue)), time.Time{})
	c.conn.WriteInt64(intValue)
}
