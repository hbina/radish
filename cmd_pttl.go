package redis

import (
	"fmt"
	"time"
)

func PttlCommand(c *Client, args [][]byte) {
	if len(args) != 2 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()
	key := string(args[1])
	db.DeleteExpired(key)

	if !db.Exists(&key) {
		c.Conn().WriteInt(-2)
		return
	}

	ttl, ok := db.Expiry(key)

	if !ok {
		c.Conn().WriteInt(-1)
		return
	}

	c.Conn().WriteInt64(int64(time.Until(ttl).Milliseconds()))
}
