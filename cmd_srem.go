package redis

import (
	"fmt"
)

// https://redis.io/commands/srem/
// SREM key member [member ...]
func SremCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError(ZeroArgumentErr)
		return
	} else if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	val, _ := c.Db().GetOrExpire(key, true)

	if val == nil {
		c.Conn().WriteInt(0)
		return
	} else if val.Type() != ValueTypeSet {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	set := val.(*Set)

	count := 0
	for i := 2; i < len(args); i++ {
		if set.Exists(string(args[i])) {
			count++
		}
		set.RemoveMember(key)
	}

	c.Conn().WriteInt(count)
}
