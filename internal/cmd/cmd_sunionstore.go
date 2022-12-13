package cmd

import (
	"fmt"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
)

// https://redis.io/commands/sunionstore/
// SUNIONSTORE destination key [key ...]
func SunionstoreCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, args[0]))
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
	union := NewSetEmpty()

	for _, key := range keys {
		maybeSet, _ := db.GetOrExpire(key, true)

		// If any of the sets are nil, then the intersections must be 0
		if maybeSet == nil {
			continue
		} else if maybeSet.Type() != types.ValueTypeSet {
			c.Conn().WriteError(pkg.WrongTypeErr)
			return
		}

		set := maybeSet.(*types.Set)

		union = union.Union(set)
	}

	db.Set(destination, union, time.Time{})

	c.Conn().WriteInt(union.Len())
}
