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

	set := val.Value().(map[string]struct{})

	count := 0
	for i := 2; i < len(args); i++ {
		_, exists := set[string(args[i])]
		if exists {
			count++
		}
		delete(set, string(args[i]))
	}

	c.Conn().WriteInt(count)
}
