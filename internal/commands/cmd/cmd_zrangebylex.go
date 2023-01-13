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

// https://redis.io/commands/zrangebylex/
// ZRANGEBYLEX key min max [LIMIT offset count]
func ZrangebylexCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 4 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	startStr := string(args[2])
	stopStr := string(args[3])

	start, startExclusive, stop, stopExclusive, notOk := util.ParseLexRange(startStr, stopStr)

	if notOk {
		c.Conn().WriteError(util.InvalidLexErr)
		return
	}

	// Parse options
	offset := 0
	limit := math.MaxInt

	for i := 4; i < len(args); i++ {
		arg := strings.ToLower(string(args[i]))
		switch arg {
		default:
			{
				c.Conn().WriteError(util.SyntaxErr)
				return
			}
		case "limit":
			{
				// Requires at least 2 more arguments
				if i+2 >= len(args) {
					c.Conn().WriteError(util.SyntaxErr)
					return
				}

				offsetStr := string(args[i+1])
				limitStr := string(args[i+2])
				i += 2

				newOffset, err := strconv.ParseInt(offsetStr, 10, 32)

				if err != nil {
					c.Conn().WriteError(util.InvalidIntErr)
					return
				}

				offset = int(newOffset)

				newLimit, err := strconv.ParseInt(limitStr, 10, 32)

				if err != nil {
					c.Conn().WriteError(util.InvalidIntErr)
					return
				}

				limit = int(newLimit)
			}
		}
	}

	maybeSet, _ := c.Db().Get(key)

	if maybeSet == nil {
		c.Conn().WriteError(util.WrongTypeErr)
		return
	}

	if maybeSet.Type() != types.ValueTypeZSet {
		c.Conn().WriteError(util.WrongTypeErr)
		return
	}

	set := maybeSet.Value().(*types.SortedSet)

	res := set.GetRangeByLex(start, stop, types.GetRangeOptions{
		Reverse:        false,
		Offset:         offset,
		Limit:          limit,
		StartExclusive: startExclusive,
		StopExclusive:  stopExclusive,
	})

	c.WriteToConn(res, false, true)
}
