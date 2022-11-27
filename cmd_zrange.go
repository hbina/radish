package redis

import (
	"fmt"
	"strconv"
	"strings"
)

// https://redis.io/commands/zrange/
// ZRANGE key start stop [BYSCORE | BYLEX] [REV] [LIMIT offset count] [WITHSCORES]
func ZrangeCommand(c *Client, args [][]byte) {
	if len(args) < 4 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	startStr := string(args[2])
	stopStr := string(args[3])

	start64, err := strconv.ParseInt(startStr, 10, 32)

	if err != nil {
		c.Conn().WriteError(InvalidIntErr)
		return
	}

	start := int(start64)

	stop64, err := strconv.ParseInt(stopStr, 10, 32)

	if err != nil {
		c.Conn().WriteError(InvalidIntErr)
		return
	}

	stop := int(stop64)

	// Parse options
	sortByLex := false
	sortByScore := false
	withScores := false
	reverse := false
	offset := 0
	limit := 0

	// TODO: Can be optimized to end when we encounter an integer
	for i := 4; i < len(args); i++ {
		arg := strings.ToLower(string(args[i]))
		switch arg {
		default:
			{
				c.Conn().WriteError(SyntaxErr)
				return
			}
		case "byscore":
			{
				if sortByLex {
					c.Conn().WriteError(SyntaxErr)
					return
				}
				sortByScore = true
			}
		case "bylex":
			{
				if sortByScore {
					c.Conn().WriteError(SyntaxErr)
					return
				}
				sortByLex = true
			}
		case "rev":
			{
				reverse = true
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

	set := maybeSet.Value().(SortedSet[string, float64, struct{}])

	res := set.GetRangeByIndex(start, stop, reverse, false)

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

	fmt.Println(reverse, offset, limit, set)
}
