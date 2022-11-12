package redis

import (
	"fmt"
	"strconv"
	"strings"
)

// https://redis.io/commands/sintercard/
// SREM key member [member ...]
// TODO: Cleanup this mess. It feels like this shouldn't be as complicated as this?
func SintercardCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError(ZeroArgumentErr)
		return
	} else if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	numberOfKeys64, err := strconv.ParseInt(string(args[1]), 10, 32)

	if err != nil {
		c.Conn().WriteError(InvalidIntErr)
		return
	}

	numberOfKeys := int(numberOfKeys64)

	if numberOfKeys < 0 {
		c.Conn().WriteError("ERR numkeys should be greater than 0")
		return
	}

	// Should not be possible to have more keys than the args passed
	if numberOfKeys > len(args)-2 {
		c.Conn().WriteError("ERR Number of keys can't be greater than number of args")
		return
	}

	// The only additional args that can be passed is LIMIT <limit>
	if numberOfKeys+2 != len(args) && numberOfKeys+4 != len(args) {
		c.Conn().WriteError(SyntaxErr)
		return
	}

	// Collect keys
	// TODO: Can optimize by removing this temporary array and use the args directly.
	keys := make([]string, 0, numberOfKeys)
	for i := 0; i < numberOfKeys; i++ {
		keys = append(keys, string(args[i+2]))
	}

	// Parse limit option
	limit := 0

	// number of keys should be equal to the length of the args minus
	// command name, number of keys, limit and limit number
	if len(args)-4 == numberOfKeys {
		limitOption := string(args[len(args)-2])
		limitValue64, err := strconv.ParseInt(string(args[len(args)-1]), 10, 32)

		if strings.ToLower(limitOption) != "limit" || err != nil || limitValue64 < 0 {
			c.Conn().WriteError(SyntaxErr)
			return
		}

		limit = int(limitValue64)
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
		return
	}

	if limit > intersection.Len() || limit == 0 {
		c.Conn().WriteInt(intersection.Len())
	} else {
		c.Conn().WriteInt(limit)
	}
}
