package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
)

// https://redis.io/commands/srem/
// SREM key member [member ...]
func SremCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()
	key := string(args[1])
	maybeSet, ttl := db.GetOrExpire(key, true)

	if maybeSet == nil {
		c.Conn().WriteInt(0)
		return
	} else if maybeSet.Type() != types.ValueTypeSet {
		c.Conn().WriteError(pkg.WrongTypeErr)
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

	c.Conn().WriteInt(count)
}
