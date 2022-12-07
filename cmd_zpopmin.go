package redis

import (
	"fmt"
	"strconv"
)

// https://redis.io/commands/zpopmin/
// ZPOPMIN key [count]
func ZpopminCommand(c *Client, args [][]byte) {
	if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])

	// Parse options
	count := 1

	if len(args) == 3 {
		countStr := string(args[2])

		count64, err := strconv.ParseInt(countStr, 10, 32)

		if err != nil {
			c.Conn().WriteError(InvalidIntErr)
			return
		}

		if count64 < 0 {
			c.Conn().WriteError(fmt.Sprintf(NegativeIntErr, "count"))
			return
		}

		count = int(count64)
	}

	db := c.Db()
	maybeSet, ttl := db.GetOrExpire(key, true)

	if maybeSet == nil {
		maybeSet = NewZSet()
	}

	if maybeSet.Type() != ValueTypeZSet {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	if count == 0 {
		c.Conn().WriteArray(0)
		return
	}

	set := maybeSet.Value().(*SortedSet)

	if count > set.Len() {
		count = set.Len()
	}

	res := set.GetRangeByRank(1, count, DefaultRangeOptions())

	for _, n := range res {
		set.Remove(n.key)
	}

	db.Set(key, NewZSetFromSs(set), ttl)

	c.Conn().WriteArray(len(res) * 2)
	for _, n := range res {
		c.Conn().WriteBulkString(n.key)
		c.Conn().WriteBulkString(fmt.Sprint(n.score))
	}
}
