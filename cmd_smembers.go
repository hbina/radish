package redis

import (
	"fmt"
)

// https://redis.io/commands/sadd/
func SmembersCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError(ZeroArgumentErr)
		return
	} else if len(args) != 2 {
		c.Conn().WriteError(fmt.Sprintf("wrong number of arguments for '%s' command", args[0]))
		return
	}

	key := string(args[1])
	maybeSet := c.Db().Get(key)

	if maybeSet == nil {
		maybeSet = NewSetEmpty()
	}

	if maybeSet.Type() != ValueTypeSet {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	set := maybeSet.Value().(map[string]struct{})

	// We already checked that there are at least 3 arguments.
	// So this should at least iterate once
	result := make([]string, 0)
	for k := range set {
		result = append(result, k)
	}

	c.Conn().WriteArray(len(result))
	for _, v := range result {
		c.Conn().WriteBulkString(v)
	}
}
