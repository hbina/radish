package redis

import (
	"fmt"
	"time"
)

// https://redis.io/commands/setnx/
// SETNX key value
// This is equivalent to calling SET key value NX
func SetNxCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
		return
	} else if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	value := string(args[2])

	db := c.Db()
	exists := db.Exists(&key)

	if exists {
		c.Conn().WriteInt(0)
		return
	}

	db.Set(key, NewString(value), time.Time{})

	c.Conn().WriteInt(1)
}
