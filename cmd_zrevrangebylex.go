package redis

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// https://redis.io/commands/zrevrangebylex/
// ZREVRANGEBYLEX key max min [LIMIT offset count]
func ZrevrangebylexCommand(c *Client, args [][]byte) {
	if len(args) < 4 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	startStr := string(args[2])
	stopStr := string(args[3])

	start, startExclusive, stop, stopExclusive, err := ParseLexRange(startStr, stopStr)

	if err {
		c.Conn().WriteError(InvalidLexErr)
		return
	}

	// Parse options
	offset := 0
	limit := math.MaxInt

	// TODO: Can be optimized to end when we encounter an integer
	for i := 4; i < len(args); i++ {
		arg := strings.ToLower(string(args[i]))
		switch arg {
		default:
			{
				c.Conn().WriteError(SyntaxErr)
				return
			}
		case "limit":
			{
				// Requires at least 2 more arguments
				if i+2 >= len(args) {
					c.Conn().WriteError(SyntaxErr)
					return
				}

				offsetStr := string(args[i+1])
				limitStr := string(args[i+2])
				i += 2

				newOffset, err := strconv.ParseInt(offsetStr, 10, 32)

				if err != nil {
					c.Conn().WriteError(InvalidIntErr)
					return
				}

				offset = int(newOffset)

				newLimit, err := strconv.ParseInt(limitStr, 10, 32)

				if err != nil {
					c.Conn().WriteError(InvalidIntErr)
					return
				}

				limit = int(newLimit)
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

	res := set.GetRangeByLex(start, stop, GetRangeOptions{
		reverse:        true,
		offset:         offset,
		limit:          limit,
		startExclusive: startExclusive,
		stopExclusive:  stopExclusive,
	})

	c.Conn().WriteArray(len(res))

	for _, ssn := range res {
		c.Conn().WriteBulkString(ssn.key)
	}
}
