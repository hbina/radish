package redis

import (
	"fmt"
	"time"
)

// https://redis.io/commands/sdiffstore/
// SDIFFSTORE destination key [key ...]
func SdiffstoreCommand(c *Client, args [][]byte) {
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
	for i := 2; i < len(args); i++ {
		keys = append(keys, string(args[i]))
	}

	db := c.Db()
	var diff *Set = nil

	for _, key := range keys {
		maybeSet, _ := db.GetOrExpire(key, true)

		// If any of the sets are nil, then the intersections must be 0
		if maybeSet == nil {
			maybeSet = NewSetEmpty()
		} else if maybeSet.Type() != ValueTypeSet {
			c.Conn().WriteError(WrongTypeErr)
			return
		}

		set := maybeSet.(*Set)

		if diff == nil {
			diff = set
		} else {
			diff = diff.Diff(set)
		}
	}

	db.Set(destination, diff, time.Time{})

	c.Conn().WriteInt(diff.Len())
}
