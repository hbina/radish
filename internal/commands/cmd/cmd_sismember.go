package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/sismember/
func SismemberCommand(c *pkg.Client, args [][]byte) {
	if len(args) != 3 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	member := string(args[2])
	maybeSet := c.Db().Get(key)

	if maybeSet == nil {
		maybeSet = types.NewSetEmpty()
	}

	if maybeSet.Type() != types.ValueTypeSet {
		c.Conn().WriteError(util.WrongTypeErr)
		return
	}

	set := maybeSet.(*types.Set)

	if set.Exists(member) {
		c.Conn().WriteInt(1)
	} else {
		c.Conn().WriteInt(0)
	}

}
