package cmd

import (
	"fmt"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/sunionstore/
// SUNIONSTORE destination key [key ...]
func SunionstoreCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 3 {
		c.WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
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
	union := types.NewSetEmpty()

	for _, key := range keys {
		maybeSet, _ := db.Get(key)

		// If any of the sets are nil, then the intersections must be 0
		if maybeSet == nil {
			continue
		} else if maybeSet.Type() != types.ValueTypeSet {
			c.WriteError(util.WrongTypeErr)
			return
		}

		set := maybeSet.(*types.Set)

		union = union.Union(set)
	}

	db.Set(destination, union, time.Time{})

	c.WriteInt(union.Len())
}
