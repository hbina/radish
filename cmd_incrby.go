package redis

import (
	"fmt"
	"strconv"
)

func IncrByCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
		return
	} else if len(args) != 3 {
		c.Conn().WriteError(fmt.Sprintf("wrong number of arguments for '%s' command", args[0]))
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
		db.Set(key, NewString(fmt.Sprintf("%d", incrBy)), nil)
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

	db.Set(key, NewString(fmt.Sprint(intValue)), nil)
	c.conn.WriteInt64(intValue)
}
