package redis

import (
	"fmt"
	"time"
)

// https://redis.io/commands/sadd/
func SaddCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError(ZeroArgumentErr)
		return
	} else if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
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
	count := 0
	for i := 2; i < len(args); i++ {
		newMember := string(args[i])
		_, found := set[newMember]
		if !found {
			set[newMember] = struct{}{}
			count++
		}

	}

	c.Db().Set(key, NewSetFromMap(set), time.Time{})

	c.Conn().WriteInt(count)
}