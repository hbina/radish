package redis

import (
	"fmt"
	"strconv"
	"time"
)

// https://redis.io/commands/decrby/
func DecrByCommand(c *Client, args [][]byte) {
	if len(args) != 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()

	key := string(args[1])
	decrBy, err := strconv.ParseInt(string(args[2]), 10, 64)

	if err != nil {
		c.Conn().WriteError(InvalidIntErr)
		return
	}

	item, exists := db.storage[key]

	if !exists {
		db.Set(key, NewString(fmt.Sprintf("%d", decrBy)), time.Time{})
		c.conn.WriteInt64(decrBy)
		return
	}

	value, ok := item.Value().(string)

	if !ok {
		c.conn.WriteError(WrongTypeErr)
		return
	}

	intValue, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		c.conn.WriteError(InvalidIntErr)
		return
	}

	intValue -= decrBy

	db.Set(key, NewString(fmt.Sprint(intValue)), time.Time{})
	c.conn.WriteInt64(intValue)
}
