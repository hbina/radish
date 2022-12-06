package redis

import (
	"fmt"
	"math"
)

// https://redis.io/commands/zremrangebyscore/
// ZREMRANGEBYSCORE key min max
func ZremrangebyscoreCommand(c *Client, args [][]byte) {
	if len(args) != 4 {
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

	db := c.Db()
	maybeSet := db.Get(key)

	if maybeSet == nil {
		maybeSet = NewZSet()
	}

	if maybeSet.Type() != ValueTypeZSet {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	set := maybeSet.Value().(*SortedSet)

	res := set.GetRangeByScore(start, stop, GetRangeOptions{
		reverse:        false,
		offset:         0,
		limit:          math.MaxInt,
		startExclusive: startExclusive,
		stopExclusive:  stopExclusive,
	})

	count := 0
	for _, r := range res {
		if set.Remove(r.key) != nil {
			count++
		}
	}

	if set.Len() == 0 {
		db.Delete(key)
	}

	c.Conn().WriteInt(count)
}
