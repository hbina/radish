package redis

import (
	"fmt"
)

// https://redis.io/commands/zcount/
// ZCOUNT key min max
func ZcountCommand(c *Client, args [][]byte) {
	if len(args) < 4 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	startStr := string(args[2])
	stopStr := string(args[3])

	start, startExclusive, stop, stopExclusive, err := ParseFloatRange(startStr, stopStr)

	if err {
		c.Conn().WriteError(InvalidFloatErr)
		return
	}

	maybeSet := c.Db().Get(key)

	if maybeSet == nil {
		maybeSet = NewZSet()
	}

	if maybeSet.Type() != ValueTypeZSet {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	set := maybeSet.Value().(*SortedSet)

	options := DefaultRangeOptions()
	options.startExclusive = startExclusive
	options.stopExclusive = stopExclusive

	res := set.GetRangeByScore(start, stop, options)

	c.Conn().WriteInt(len(res))
}
