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
		c.WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	startStr := string(args[2])
	stopStr := string(args[3])

	start, startExclusive, stop, stopExclusive, err := util.ParseFloatRange(startStr, stopStr)

	if err {
		c.WriteError(util.InvalidFloatErr)
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
				c.WriteError(util.SyntaxErr)
				return
			}
		case "limit":
			{
				// Requires at least 2 more arguments
				if i+2 >= len(args) {
					c.WriteError(util.SyntaxErr)
					return
				}

				offsetStr := string(args[i+1])
				limitStr := string(args[i+2])
				i += 2

				newOffset, err := strconv.ParseInt(offsetStr, 10, 32)

				if err != nil {
					c.WriteError(util.InvalidIntErr)
					return
				}

				offset = int(newOffset)

				newLimit, err := strconv.ParseInt(limitStr, 10, 32)

				if err != nil {
					c.WriteError(util.InvalidIntErr)
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

	maybeSet, _ := c.Db().Get(key)

	if maybeSet == nil {
		maybeSet = types.NewZSet()
	}

	if maybeSet.Type() != types.ValueTypeZSet {
		c.WriteError(util.WrongTypeErr)
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
		c.WriteArray(len(res) * 2)

		for _, ssn := range res {
			c.WriteBulkString(ssn.Key)
			c.WriteBulkString(fmt.Sprint(ssn.Score))
		}
	} else {
		c.WriteArray(len(res))

		for _, ssn := range res {
			c.WriteBulkString(ssn.Key)
		}
	}
}
