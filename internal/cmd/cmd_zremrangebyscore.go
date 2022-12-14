package cmd

import (
	"fmt"
	"math"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/zremrangebyscore/
// ZREMRANGEBYSCORE key min max
func ZremrangebyscoreCommand(c *pkg.Client, args [][]byte) {
	if len(args) != 4 {
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	startStr := string(args[2])
	stopStr := string(args[3])

	start, startExclusive, stop, stopExclusive, err := util.ParseFloatRange(startStr, stopStr)

	if err {
		c.Conn().WriteError(pkg.InvalidFloatErr)
		return
	}

	db := c.Db()
	maybeSet := db.Get(key)

	if maybeSet == nil {
		maybeSet = types.NewZSet()
	}

	if maybeSet.Type() != types.ValueTypeZSet {
		c.Conn().WriteError(pkg.WrongTypeErr)
		return
	}

	set := maybeSet.Value().(*types.SortedSet)

	res := set.GetRangeByScore(start, stop, types.GetRangeOptions{
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
