package redis

import (
	"fmt"
	"strconv"
	"time"
)

// https://redis.io/commands/incrby/
func IncrByCommand(c *Client, args [][]byte) {
	if len(args) != 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()

	key := string(args[1])
	incrBy, err := strconv.ParseInt(string(args[2]), 10, 64)

	if err != nil {
		c.Conn().WriteError(InvalidIntErr)
		return
	}

	item, exists := db.storage[key]

	if !exists {
		db.Set(key, NewString(fmt.Sprintf("%d", incrBy)), time.Time{})
		c.conn.WriteInt64(incrBy)
		return
	}

	value, ok := item.Value().(*string)

	if !ok {
		c.conn.WriteError(WrongTypeErr)
		return
	}

	intValue, err := strconv.ParseInt(*value, 10, 64)

	if err != nil {
		c.conn.WriteError(InvalidIntErr)
		return
	}

	intValue += incrBy

	db.Set(key, NewString(fmt.Sprint(intValue)), time.Time{})
	c.conn.WriteInt64(intValue)
}
