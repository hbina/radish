package redis

import (
	"fmt"
)

// https://redis.io/commands/zrem/
// ZREM key member [member ...]
func ZremCommand(c *Client, args [][]byte) {
	if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	db := c.Db()

	maybeSet, ttl := db.GetOrExpire(key, true)

	if maybeSet == nil {
		maybeSet = NewZSet()
	}

	if maybeSet.Type() != ValueTypeZSet {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	set := maybeSet.(*ZSet)

	count := 0
	for i := 2; i < len(args); i++ {
		res := set.inner.Remove(string(args[i]))
		if res != nil {
			count++
		}
	}

	db.Set(key, set, ttl)

	c.Conn().WriteInt(count)
}
