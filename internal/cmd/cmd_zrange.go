package cmd

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
)

// https://redis.io/commands/zrange/
// ZRANGE key start stop [BYSCORE | BYLEX] [REV] [LIMIT offset count] [WITHSCORES]
func ZrangeCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 4 {
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	startStr := string(args[2])
	stopStr := string(args[3])

	// Parse options
	sortByLex := false
	sortByScore := false
	withScores := false
	reverse := false
	offset := 0
	limit := math.MaxInt

	for i := 4; i < len(args); i++ {
		arg := strings.ToLower(string(args[i]))

		switch arg {
		default:
			{
				c.Conn().WriteError(pkg.SyntaxErr)
				return
			}
		case "byscore":
			{
				if sortByLex {
					c.Conn().WriteError(pkg.SyntaxErr)
					return
				}
				sortByScore = true
			}
		case "bylex":
			{
				if sortByScore {
					c.Conn().WriteError(pkg.SyntaxErr)
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
					c.Conn().WriteError(pkg.SyntaxErr)
					return
				}

				offsetStr := string(args[i+1])
				limitStr := string(args[i+2])
				i += 2

				newOffset, err := strconv.ParseInt(offsetStr, 10, 32)

				if err != nil {
					c.Conn().WriteError(pkg.InvalidIntErr)
					return
				}

				offset = int(newOffset)

				newLimit, err := strconv.ParseInt(limitStr, 10, 32)

				if err != nil {
					c.Conn().WriteError(pkg.InvalidIntErr)
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

	if maybeSet.Type() != types.ValueTypeZSet {
		c.Conn().WriteError(pkg.WrongTypeErr)
		return
	}

	set := maybeSet.Value().(*types.SortedSet)

	var res []*types.SortedSetNode

	if sortByLex {
		start, startExclusive, stop, stopExclusive, err := ParseLexRange(startStr, stopStr)

		if err {
			c.Conn().WriteError(InvalidLexErr)
			return
		}

		res = set.GetRangeByLex(start, stop, GetRangeOptions{
			reverse:        reverse,
			offset:         offset,
			limit:          limit,
			startExclusive: startExclusive,
			stopExclusive:  stopExclusive,
		})
	} else if sortByScore {
		start, startExclusive, stop, stopExclusive, err := ParseFloatRange(startStr, stopStr)

		if err {
			c.Conn().WriteError(pkg.InvalidFloatErr)
			return
		}

		res = set.GetRangeByScore(start, stop, GetRangeOptions{
			reverse:        reverse,
			offset:         offset,
			limit:          limit,
			startExclusive: startExclusive,
			stopExclusive:  stopExclusive,
		})
	} else {
		start, startExclusive, stop, stopExclusive, err := ParseIntRange(startStr, stopStr)

		if err {
			c.Conn().WriteError(pkg.InvalidIntErr)
			return
		}

		res = set.GetRangeByIndex(start, stop, GetRangeOptions{
			reverse:        reverse,
			offset:         offset,
			limit:          limit,
			startExclusive: startExclusive,
			stopExclusive:  stopExclusive,
		})
	}

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
