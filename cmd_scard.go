package redis

import (
	"fmt"
)

// https://redis.io/commands/scard/
// SCARD key
func ScardCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError(ZeroArgumentErr)
		return
	} else if len(args) != 2 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])

	db := c.Db()

	maybeSet, _ := db.GetOrExpire(key, true)

	if maybeSet == nil {
		c.Conn().WriteInt(0)
		return
	} else if maybeSet.Type() != ValueTypeSet {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	set := maybeSet.(*Set)

	c.Conn().WriteInt(set.Len())
}
