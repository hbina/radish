package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/zrem/
// ZREM key member [member ...]
func ZremCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 3 {
		c.WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	db := c.Db()

	maybeSet, ttl := db.Get(key)

	if maybeSet == nil {
		maybeSet = types.NewZSet()
	}

	if maybeSet.Type() != types.ValueTypeZSet {
		c.WriteError(util.WrongTypeErr)
		return
	}

	set := maybeSet.(*types.ZSet)

	count := 0
	for i := 2; i < len(args); i++ {
		res := set.Remove(string(args[i]))
		if res != nil {
			count++
		}
	}

	db.Set(key, set, ttl)

	c.WriteInt(count)
}
