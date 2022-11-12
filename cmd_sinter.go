package redis

import (
	"fmt"
)

// https://redis.io/commands/sinter/
// SINTER key [key ...]
func SinterCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError(ZeroArgumentErr)
		return
	} else if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	// Collect keys
	// TODO: Can optimize by removing this temporary array and use the args directly.
	keys := make([]string, 0, len(args)-1)
	for i := 1; i < len(args); i++ {
		keys = append(keys, string(args[i]))
	}

	intersection := genericSinterCommand(c, keys)

	if intersection == nil {
		return
	}

	c.Conn().WriteArray(intersection.Len())
	intersection.ForEachF(func(a string) {
		c.Conn().WriteBulkString(a)
	})
}

// https://redis.io/commands/sinter/
// SINTER key [key ...]
func genericSinterCommand(c *Client, keys []string) *Set {
	db := c.Db()

	result := NewSetEmpty()

	// TODO: Is it possible to optimize using the fact that we know what the
	// upper bound is?
	for i, key := range keys {
		maybeSet, _ := db.GetOrExpire(key, true)

		// If any of the sets are nil, then the intersections must be 0
		if maybeSet == nil {
			c.Conn().WriteInt(0)
			return nil
		} else if maybeSet.Type() != ValueTypeSet {
			c.Conn().WriteError(WrongTypeErr)
			return nil
		}

		set := maybeSet.(*Set)

		if i == 0 {
			result = set
		} else {
			result = result.Intersect(set)
		}

		// TODO: Optimization to return nil early by checking if intersection is empty
	}

	return result
}
