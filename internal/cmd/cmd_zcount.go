package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/zcount/
// ZCOUNT key min max
func ZcountCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 4 {
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

	maybeSet := c.Db().Get(key)

	if maybeSet == nil {
		maybeSet = types.NewZSet()
	}

	if maybeSet.Type() != types.ValueTypeZSet {
		c.Conn().WriteError(pkg.WrongTypeErr)
		return
	}

	set := maybeSet.Value().(*types.SortedSet)

	options := types.DefaultRangeOptions()
	options.StartExclusive = startExclusive
	options.StopExclusive = stopExclusive

	res := set.GetRangeByScore(start, stop, options)

	c.Conn().WriteInt(len(res))
}