package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/sdiff/
// SDIFF key [key ...]
func SdiffCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	// Collect keys
	// TODO: Can optimize by removing this temporary array and use the args directly.
	keys := make([]string, 0, len(args)-2)
	for i := 1; i < len(args); i++ {
		keys = append(keys, string(args[i]))
	}

	db := c.Db()
	var diff *types.Set = nil

	// TODO: Is it possible to optimize using the fact that we know what the
	// upper bound is?
	for _, key := range keys {

		maybeSet, _ := db.Get(key)

		// If any of the sets are nil, then the intersections must be 0
		if maybeSet == nil {
			maybeSet = types.NewSetEmpty()
		} else if maybeSet.Type() != types.ValueTypeSet {
			c.Conn().WriteError(util.WrongTypeErr)
			return
		}

		set := maybeSet.(*types.Set)

		if diff == nil {
			diff = set
		} else {
			diff = diff.Diff(set)
		}
	}

	c.Conn().WriteArray(diff.Len())
	diff.ForEachF(func(a string) bool {
		c.Conn().WriteBulkString(a)
		return true
	})
}
