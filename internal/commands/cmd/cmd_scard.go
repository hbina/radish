package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/scard/
// SCARD key
func ScardCommand(c *pkg.Client, args [][]byte) {
	if len(args) != 2 {
		c.WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])

	db := c.Db()

	maybeSet, _ := db.Get(key)

	if maybeSet == nil {
		c.WriteInt(0)
		return
	} else if maybeSet.Type() != types.ValueTypeSet {
		c.WriteError(util.WrongTypeErr)
		return
	}

	set := maybeSet.(*types.Set)

	c.WriteInt(set.Len())
}
