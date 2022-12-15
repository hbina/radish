package cmd

import (
	"fmt"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/sinterstore/
// SREM key member [member ...]
// TODO: Cleanup this mess. It feels like this shouldn't be as complicated as this?
func SinterstoreCommand(c *pkg.Client, args [][]byte) {
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
	var intersection *types.Set = nil

	//NOTE: Cannot optimize the following loop because we need to verify that each keys consist of sets/empty.
	for _, key := range keys {
		maybeSet, _ := db.GetOrExpire(key, true)

		// If any of the sets are nil, then the intersections must be 0
		if maybeSet == nil {
			maybeSet = types.NewSetEmpty()
		} else if maybeSet.Type() != types.ValueTypeSet {
			c.Conn().WriteError(util.WrongTypeErr)
			return
		}

		set := maybeSet.(*types.Set)

		if intersection == nil {
			intersection = set
		} else {
			intersection = intersection.Intersect(set)
		}

	}

	if intersection == nil || intersection.Len() == 0 {
		// This should not be possible but just to make it look nicer.
		db.Delete(destination)
		c.Conn().WriteInt(0)
		return
	} else {
		db.Set(destination, intersection, time.Time{})
	}

	c.Conn().WriteInt(intersection.Len())
}
