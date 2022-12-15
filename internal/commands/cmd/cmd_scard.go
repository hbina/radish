package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
)

// https://redis.io/commands/scard/
// SCARD key
func ScardCommand(c *pkg.Client, args [][]byte) {
	if len(args) != 2 {
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])

	db := c.Db()

	maybeSet, _ := db.GetOrExpire(key, true)

	if maybeSet == nil {
		c.Conn().WriteInt(0)
		return
	} else if maybeSet.Type() != types.ValueTypeSet {
		c.Conn().WriteError(pkg.WrongTypeErr)
		return
	}

	set := maybeSet.(*types.Set)

	c.Conn().WriteInt(set.Len())
}
