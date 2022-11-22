package redis

import (
	"fmt"
	"strconv"
	"time"
)

// // https://redis.io/commands/decrbyfloat/
func DecrByFloatCommand(c *Client, args [][]byte) {
	if len(args) != 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()

	key := string(args[1])
	decrBy, err := strconv.ParseFloat(string(args[2]), 64)

	if err != nil {
		c.Conn().WriteError(InvalidFloatErr)
		return
	}

	item, exists := db.storage[key]

	if !exists {
		decrByStr := strconv.FormatFloat(decrBy, 'f', -1, 64)
		db.Set(key, NewString(decrByStr), time.Time{})
		c.conn.WriteString(fmt.Sprintf("\"%s\"", decrByStr))
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

	floatValue -= decrBy

	floatValueStr := strconv.FormatFloat(floatValue, 'f', -1, 64)
	db.Set(key, NewString(floatValueStr), time.Time{})
	c.conn.WriteString(fmt.Sprintf("\"%s\"", floatValueStr))
}
