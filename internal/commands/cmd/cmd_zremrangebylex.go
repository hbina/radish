package cmd

import (
	"fmt"
	"math"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/zremrangebylex/
// ZREMRANGEBYLEX key min max
func ZremrangebylexCommand(c *pkg.Client, args [][]byte) {
	if len(args) != 4 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	startStr := string(args[2])
	stopStr := string(args[3])

	start, startExclusive, stop, stopExclusive, err := util.ParseLexRange(startStr, stopStr)

	if err {
		c.Conn().WriteError(util.InvalidFloatErr)
		return
	}

	db := c.Db()
	maybeSet := db.Get(key)

	if maybeSet == nil {
		maybeSet = types.NewZSet()
	}

	if maybeSet.Type() != types.ValueTypeZSet {
		c.Conn().WriteError(util.WrongTypeErr)
		return
	}

	set := maybeSet.Value().(*types.SortedSet)

	res := set.GetRangeByLex(start, stop, types.GetRangeOptions{
		Reverse:        false,
		Offset:         0,
		Limit:          math.MaxInt,
		StartExclusive: startExclusive,
		StopExclusive:  stopExclusive,
	})

	count := 0
	for _, r := range res {
		if set.Remove(r.Key) != nil {
			count++
		}
	}

	if set.Len() == 0 {
		db.Delete(key)
	}

	c.Conn().WriteInt(count)
}
