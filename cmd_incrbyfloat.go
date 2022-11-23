package redis

import (
	"fmt"
	"strconv"
	"time"
)

// https://redis.io/commands/incrbyfloat/
func IncrByFloatCommand(c *Client, args [][]byte) {
	if len(args) != 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()

	key := string(args[1])
	incrBy, err := strconv.ParseFloat(string(args[2]), 64)

	if err != nil {
		c.Conn().WriteError(InvalidFloatErr)
		return
	}

	item, exists := db.storage[key]

	if !exists {
		incrByStr := strconv.FormatFloat(incrBy, 'f', -1, 64)
		db.Set(key, NewString(incrByStr), time.Time{})
		c.conn.WriteString(fmt.Sprintf("\"%s\"", incrByStr))
		return
	}

	value, ok := item.Value().(*string)

	if !ok {
		c.conn.WriteError(WrongTypeErr)
		return
	}

	floatValue, err := strconv.ParseFloat(*value, 64)

	if err != nil {
		c.conn.WriteError(InvalidFloatErr)
		return
	}

	floatValue += incrBy

	floatValueStr := strconv.FormatFloat(floatValue, 'f', -1, 64)
	db.Set(key, NewString(floatValueStr), time.Time{})
	c.conn.WriteString(fmt.Sprintf("\"%s\"", floatValueStr))
}
