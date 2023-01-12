package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/srem/
// SREM key member [member ...]
func SremCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 3 {
		c.WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()
	key := string(args[1])
	maybeSet, ttl := db.Get(key)

	if maybeSet == nil {
		c.WriteInt(0)
		return
	} else if maybeSet.Type() != types.ValueTypeSet {
		c.WriteError(util.WrongTypeErr)
		return
	}

	set := maybeSet.(*types.Set)

	count := 0
	for i := 2; i < len(args); i++ {
		if set.Exists(string(args[i])) {
			count++
		}
		set.RemoveMember(string(args[i]))
	}

	db.Set(key, set, ttl)

	c.WriteInt(count)
}
