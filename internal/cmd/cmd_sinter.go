package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
)

// https://redis.io/commands/sinter/
// SINTER key [key ...]
func SinterCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, args[0]))
		return
	}

	// Collect keys
	// TODO: Can optimize by removing this temporary array and use the args directly.
	keys := make([]string, 0, len(args)-1)
	for i := 1; i < len(args); i++ {
		keys = append(keys, string(args[i]))
	}

	db := c.Db()

	var intersection *types.Set = nil

	for _, key := range keys {
		maybeSet, _ := db.GetOrExpire(key, true)

		// If any of the sets are nil, then the intersections must be 0
		if maybeSet == nil {
			maybeSet = types.NewSetEmpty()
		} else if maybeSet.Type() != types.ValueTypeSet {
			c.Conn().WriteError(pkg.WrongTypeErr)
			return
		}

		set := maybeSet.(*types.Set)

		if intersection == nil {
			intersection = set
		} else {
			intersection = intersection.Intersect(set)
		}
	}

	if intersection == nil {
		c.Conn().WriteArray(0)
		return
	}

	c.Conn().WriteArray(intersection.Len())
	intersection.ForEachF(func(a string) {
		c.Conn().WriteBulkString(a)
	})
}
