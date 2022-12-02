package redis

import (
	"fmt"
	"strings"
)

// https://redis.io/commands/zrevrange/
// ZREVRANGE key start stop [WITHSCORES]
func ZrevrangeCommand(c *Client, args [][]byte) {
	if len(args) < 4 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	startStr := string(args[2])
	stopStr := string(args[3])

	start, startExclusive, stop, stopExclusive, err := ParseIntRange(startStr, stopStr)

	if err {
		c.Conn().WriteError(InvalidIntErr)
		return
	}

	// Parse options
	withScores := false

	// TODO: Can be optimized to end when we encounter an integer
	for i := 4; i < len(args); i++ {
		arg := strings.ToLower(string(args[i]))
		switch arg {
		default:
			{
				c.Conn().WriteError(SyntaxErr)
				return
			}
		case "withscores":
			{
				withScores = true
			}
		}
	}

	maybeSet := c.Db().Get(key)

	if maybeSet == nil {
		maybeSet = NewZSet()
	}

	if maybeSet.Type() != ValueTypeZSet {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	set := maybeSet.Value().(SortedSet)

	options := DefaultRangeOptions()
	options.reverse = true
	options.startExclusive = startExclusive
	options.stopExclusive = stopExclusive

	res := set.GetRangeByIndex(start, stop, options)

	if withScores {
		c.Conn().WriteArray(len(res) * 2)

		for _, ssn := range res {
			c.Conn().WriteBulkString(ssn.key)
			c.Conn().WriteBulkString(fmt.Sprint(ssn.score))
		}
	} else {
		c.Conn().WriteArray(len(res))

		for _, ssn := range res {
			c.Conn().WriteBulkString(ssn.key)
		}
	}
}
