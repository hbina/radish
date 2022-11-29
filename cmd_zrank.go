package redis

import (
	"fmt"
)

// https://redis.io/commands/zrank/
// ZRANK key member
func ZrankCommand(c *Client, args [][]byte) {
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

	memberRank := set.inner.FindRankOfKey(memberKey)

	if memberRank == 0 {
		c.Conn().WriteNull()
		return
	}

	c.Conn().WriteString(fmt.Sprint(memberRank - 1))
}
