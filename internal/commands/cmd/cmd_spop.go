package cmd

import (
	"fmt"
	"strconv"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/spop/
// SPOP key [count]
func SpopCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	var count *int = nil

	if len(args) == 3 {
		count64, err := strconv.ParseInt(string(args[2]), 10, 32)

		if err != nil || count64 < 0 {
			c.Conn().WriteError(util.InvalidIntErr)
		}

		count32 := int(count64)
		count = &count32
	}

	db := c.Db()

	maybeSet, _ := db.GetOrExpire(key, true)

	// If any of the sets are nil, then the intersections must be 0
	if maybeSet == nil {
		c.Conn().WriteNull()
		return
	} else if maybeSet.Type() != types.ValueTypeSet {
		c.Conn().WriteError(util.WrongTypeErr)
		return
	}

	set := maybeSet.(*types.Set)

	// TODO: If count larger than set, then just delete set
	if count != nil {
		removed := make([]string, 0)

		for i := 0; i < *count; i++ {
			member := set.Pop()

			if member != nil {
				removed = append(removed, *member)
			}
		}

		c.Conn().WriteArray(len(removed))
		for _, k := range removed {
			c.Conn().WriteBulkString(k)
		}
	} else {
		member := set.Pop()

		if member != nil {
			c.Conn().WriteBulkString(*member)
		} else {
			c.Conn().WriteNull()
		}
	}
}
