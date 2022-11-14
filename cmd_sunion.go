package redis

import (
	"fmt"
)

// https://redis.io/commands/sunion/
// SUNION key [key ...]
func SunionCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError(ZeroArgumentErr)
		return
	} else if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	// Collect keys
	// TODO: Can optimize by removing this temporary array and use the args directly.
	keys := make([]string, 0, len(args)-2)
	for i := 1; i < len(args); i++ {
		keys = append(keys, string(args[i]))
	}

	db := c.Db()
	union := NewSetEmpty()

	for _, key := range keys {
		maybeSet, _ := db.GetOrExpire(key, true)

		// If any of the sets are nil, then the intersections must be 0
		if maybeSet == nil {
			continue
		} else if maybeSet.Type() != ValueTypeSet {
			c.Conn().WriteError(WrongTypeErr)
			return
		}

		set := maybeSet.(*Set)
		union = union.Union(set)
	}

	c.Conn().WriteArray(union.Len())
	union.ForEachF(func(a string) {
		c.Conn().WriteBulkString(a)
	})
}
