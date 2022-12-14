package cmd

import (
	"fmt"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/sdiffstore/
// SDIFFSTORE destination key [key ...]
func SdiffstoreCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
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
	var diff *types.Set = nil

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

	db.Set(destination, diff, time.Time{})

	c.Conn().WriteInt(diff.Len())
}
