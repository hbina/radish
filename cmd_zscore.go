package redis

import (
	"fmt"
	"math"
)

// https://redis.io/commands/zscore/
// ZSCORE key member
func ZscoreCommand(c *Client, args [][]byte) {
	if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	memberKey := string(args[2])
	maybeSet := c.Db().Get(key)

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
		c.Conn().WriteNull()
		return
	}

	if math.IsNaN(maybeMember.score) {
		c.Conn().WriteString("nan")
	} else if math.IsInf(maybeMember.score, -1) {
		c.Conn().WriteString("-inf")
	} else if math.IsInf(maybeMember.score, 1) {
		c.Conn().WriteString("inf")
	} else {
		c.Conn().WriteString(fmt.Sprint(maybeMember.score))
	}
}
