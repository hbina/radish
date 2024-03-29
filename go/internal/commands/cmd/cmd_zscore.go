package cmd

import (
	"fmt"
	"math"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/zscore/
// ZSCORE key member
func ZscoreCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	memberKey := string(args[2])
	maybeSet, _ := c.Db().Get(key)

	if maybeSet == nil {
		maybeSet = types.NewZSet()
	}

	if maybeSet.Type() != types.ValueTypeZSet {
		c.Conn().WriteError(util.WrongTypeErr)
		return
	}

	set := maybeSet.(*types.ZSet)

	maybeMember := set.GetByKey(memberKey)

	if maybeMember == nil {
		if c.R3 {
			c.Conn().WriteNull()
		} else {
			c.Conn().WriteNullBulk()
		}
		return
	}

	if math.IsNaN(maybeMember.Score) {
		c.Conn().WriteString("nan")
	} else if math.IsInf(maybeMember.Score, -1) {
		c.Conn().WriteString("-inf")
	} else if math.IsInf(maybeMember.Score, 1) {
		c.Conn().WriteString("inf")
	} else {
		c.Conn().WriteString(fmt.Sprint(maybeMember.Score))
	}
}
