package redis

import (
	"fmt"
)

// https://redis.io/commands/sismember/
func SismemberCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError(ZeroArgumentErr)
		return
	} else if len(args) != 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	member := string(args[2])
	maybeSet := c.Db().Get(key)

	if maybeSet == nil {
		maybeSet = NewSetEmpty()
	}

	if maybeSet.Type() != ValueTypeSet {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	set := maybeSet.(*Set)

	if set.Exists(member) {
		c.Conn().WriteInt(1)
	} else {
		c.Conn().WriteInt(0)
	}

}
