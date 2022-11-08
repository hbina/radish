package redis

import (
	"fmt"
	"time"
)

func TtlCommand(c *Client, args [][]byte) {
	if len(args) != 2 {
		c.Conn().WriteError(fmt.Sprintf("wrong number of arguments (given %d, expected 1)", len(args)-1))
		return
	}

	db := c.Db()
	key := string(args[1])
	db.DeleteExpired(&key)
	if !db.Exists(&key) {
		c.Conn().WriteInt(-2)
		return
	}

	t := db.Expiry(&key)
	if t.IsZero() {
		c.Conn().WriteInt(-1)
		return
	}

	c.Conn().WriteInt64(int64(time.Until(t).Seconds()))
}
