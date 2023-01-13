package cmd

import (
	"fmt"
	"strconv"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/zpopmin/
// ZPOPMIN key [count]
func ZpopmaxCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])

	// Parse options
	count := 1
	countSet := false

	if len(args) == 3 {
		countStr := string(args[2])

		count64, err := strconv.ParseInt(countStr, 10, 32)

		if err != nil {
			c.Conn().WriteError(util.InvalidIntErr)
			return
		}

		if count64 < 0 {
			c.Conn().WriteError(fmt.Sprintf(util.MustBePositiveErr, "count"))
			return
		}

		count = int(count64)
		countSet = true
	}

	db := c.Db()
	maybeSet, ttl := db.Get(key)

	if maybeSet == nil {
		maybeSet = types.NewZSet()
	}

	if maybeSet.Type() != types.ValueTypeZSet {
		c.Conn().WriteError(util.WrongTypeErr)
		return
	}

	if count == 0 {
		c.Conn().WriteArray(0)
		return
	}

	set := maybeSet.Value().(*types.SortedSet)

	if count > set.Len() {
		count = set.Len()
	}

	options := types.DefaultRangeOptions()
	options.Reverse = true
	res := set.GetRangeByRank(set.Len()+1-count, set.Len(), options)

	for _, n := range res {
		set.Remove(n.Key)
	}

	db.Set(key, types.NewZSetFromSs(set), ttl)

	if len(res) == 0 {
		c.Conn().WriteArray(0)
	} else if !countSet && c.R3 {
		c.Conn().WriteArray(2)
		c.Conn().WriteBulkString(res[0].Key)
		c.Conn().WriteFloat64(res[0].Score)
	} else {
		c.WriteToConn(res, true)
	}
}
