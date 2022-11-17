package redis

import (
	"fmt"
	"time"
)

// https://redis.io/commands/ttl/
// TTL key
func TtlCommand(c *Client, args [][]byte) {
	if len(args) != 2 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()
	key := string(args[1])

	item, ttl := db.GetOrExpire(key, true)

	if item != nil {
		c.Conn().WriteInt(-2)
		return
	} else if item != nil && time.Time.IsZero(ttl) {
		c.Conn().WriteInt(-1)
		return
	}

	c.Conn().WriteInt64(int64(time.Until(ttl).Seconds()))
}
