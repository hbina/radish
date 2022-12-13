package cmd

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
)

// https://redis.io/commands/zrevrangebylex/
// ZREVRANGEBYLEX key max min [LIMIT offset count]
func ZrevrangebylexCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 4 {
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, args[0]))
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
