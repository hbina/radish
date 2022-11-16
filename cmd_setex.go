package redis

import (
	"fmt"
	"time"
)

// https://redis.io/commands/setex/
// SETEX key seconds value
// This is equivalent to calling `SET key value EX seconds`
func SetexCommand(c *Client, args [][]byte) {
	if len(args) != 4 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	seconds := string(args[2])
	value := string(args[3])

	newTtl, err := ParseTtlFromUnitTime(seconds, int64(time.Second))

	if err != nil {
		c.Conn().WriteError(InvalidIntErr)
		return
	}

	db := c.Db()

	db.Set(key, NewString(value), newTtl)

	c.Conn().WriteString("OK")
}
