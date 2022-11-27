package redis

import (
	"fmt"
	"math"
	"strconv"
)

// https://redis.io/commands/zincrby/
// ZINCRBY key increment member
func ZincrbyCommand(c *Client, args [][]byte) {
	if len(args) != 4 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	incrementStr := string(args[2])
	memberKey := string(args[3])
	db := c.Db()

	increment, err := strconv.ParseFloat(incrementStr, 64)

	if err != nil || math.IsNaN(increment) {
		c.Conn().WriteError(InvalidFloatErr)
		return
	}

	maybeSet, ttl := db.GetOrExpire(key, true)

	if maybeSet == nil {
		maybeSet = NewZSet()
	}

	if maybeSet.Type() != ValueTypeZSet {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	set := maybeSet.(*ZSet)

	maybeMember := set.inner.GetByKey(memberKey)

	if maybeMember == nil {
		set.inner.AddOrUpdate(memberKey, increment, struct{}{})
		db.Set(key, set, ttl)
		c.Conn().WriteString(fmt.Sprint(increment))
	} else {
		newScore := maybeMember.Score() + increment

		if math.IsNaN(newScore) {
			c.Conn().WriteError("ERR resulting score is not a number (NaN)")
			return
		}
		set.inner.AddOrUpdate(memberKey, newScore, struct{}{})
		db.Set(key, set, ttl)
		c.Conn().WriteString(fmt.Sprint(newScore))
	}
}
