package redis

import (
	"fmt"
	"time"
)

// https://redis.io/commands/sinterstore/
// SREM key member [member ...]
// TODO: Cleanup this mess. It feels like this shouldn't be as complicated as this?
func SinterstoreCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError(ZeroArgumentErr)
		return
	} else if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	destination := string(args[1])

	// Collect keys
	// TODO: Can optimize by removing this temporary array and use the args directly.
	keys := make([]string, 0, len(args)-2)
	for i := 2; i < len(args)-2; i++ {
		keys = append(keys, string(args[i+2]))
	}

	db := c.Db()

	intersection := NewSetEmpty()

	// TODO: Is it possible to optimize using the fact that we know what the
	// upper bound is?
	for i, key := range keys {
		maybeSet, _ := db.GetOrExpire(key, true)

		// If any of the sets are nil, then the intersections must be 0
		if maybeSet == nil {
			c.Conn().WriteInt(0)
			return
		} else if maybeSet.Type() != ValueTypeSet {
			c.Conn().WriteError(WrongTypeErr)
			return
		}

		set := maybeSet.(*Set)

		if i == 0 {
			intersection = set
		} else {
			intersection = intersection.Intersect(set)
		}

		// TODO: Optimization to return nil early by checking if intersection is empty
	}

	if intersection == nil {
		db.Set(destination, NewSetEmpty(), time.Time{})
	} else {
		db.Set(destination, intersection, time.Time{})
	}

	c.Conn().WriteInt(intersection.Len())
}
