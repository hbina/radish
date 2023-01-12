package cmd

import (
	"fmt"
	"math"
	"strconv"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/zincrby/
// ZINCRBY key increment member
func ZincrbyCommand(c *pkg.Client, args [][]byte) {
	if len(args) != 4 {
		c.WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	incrementStr := string(args[2])
	memberKey := string(args[3])
	db := c.Db()

	increment, err := strconv.ParseFloat(incrementStr, 64)

	if err != nil || math.IsNaN(increment) {
		c.WriteError(util.InvalidFloatErr)
		return
	}

	maybeSet, ttl := db.Get(key)

	if maybeSet == nil {
		maybeSet = types.NewZSet()
	}

	if maybeSet.Type() != types.ValueTypeZSet {
		c.WriteError(util.WrongTypeErr)
		return
	}

	set := maybeSet.(*types.ZSet)

	maybeMember := set.GetByKey(memberKey)

	if maybeMember == nil {
		set.AddOrUpdate(memberKey, increment)
		db.Set(key, set, ttl)
		c.WriteString(fmt.Sprint(increment))
	} else {
		newScore := maybeMember.Score + increment

		if math.IsNaN(newScore) {
			c.WriteError("ERR resulting score is not a number (NaN)")
			return
		}
		set.AddOrUpdate(memberKey, newScore)
		db.Set(key, set, ttl)
		c.WriteString(fmt.Sprint(newScore))
	}
}
