package cmd

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/zrevrangebyscore/
// ZREVRANGEBYSCORE key min max [WITHSCORES] [LIMIT offset count]
func ZrevrangebyscoreCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 4 {
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	startStr := string(args[2])
	stopStr := string(args[3])

	start, startExclusive, stop, stopExclusive, err := util.ParseFloatRange(startStr, stopStr)

	if err {
		c.Conn().WriteError(pkg.InvalidFloatErr)
		return
	}

	// Parse options
	withScores := false
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
		maybeSet = types.NewZSet()
	}

	if maybeSet.Type() != types.ValueTypeZSet {
		c.Conn().WriteError(pkg.WrongTypeErr)
		return
	}

	set := maybeSet.Value().(*types.SortedSet)

	res := set.GetRangeByScore(start, stop, types.GetRangeOptions{
		Reverse:        true,
		Offset:         offset,
		Limit:          limit,
		StartExclusive: startExclusive,
		StopExclusive:  stopExclusive,
	})

	if withScores {
		c.Conn().WriteArray(len(res) * 2)

		for _, ssn := range res {
			c.Conn().WriteBulkString(ssn.Key)
			c.Conn().WriteBulkString(fmt.Sprint(ssn.Score))
		}
	} else {
		c.Conn().WriteArray(len(res))

		for _, ssn := range res {
			c.Conn().WriteBulkString(ssn.Key)
		}
	}
}
