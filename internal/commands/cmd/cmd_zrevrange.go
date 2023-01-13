package cmd

import (
	"fmt"
	"strings"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/zrevrange/
// ZREVRANGE key start stop [WITHSCORES]
func ZrevrangeCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 4 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	startStr := string(args[2])
	stopStr := string(args[3])

	start, startExclusive, stop, stopExclusive, notOk := util.ParseIntRange(startStr, stopStr)

	if notOk {
		c.Conn().WriteError(util.InvalidIntErr)
		return
	}

	// Parse options
	withScores := false

	for i := 4; i < len(args); i++ {
		arg := strings.ToLower(string(args[i]))
		switch arg {
		default:
			{
				c.Conn().WriteError(util.SyntaxErr)
				return
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
		c.Conn().WriteError(util.WrongTypeErr)
		return
	}

	set := maybeSet.Value().(*types.SortedSet)

	options := types.DefaultRangeOptions()
	options.Reverse = true
	options.StartExclusive = startExclusive
	options.StopExclusive = stopExclusive

	res := set.GetRangeByIndex(start, stop, options)
	c.WriteToConn(res, withScores, true)
}
