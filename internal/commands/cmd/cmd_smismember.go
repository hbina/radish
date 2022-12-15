package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
)

// https://redis.io/commands/smismember/
func SmismemberCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	maybeSet := c.Db().Get(key)

	if maybeSet == nil {
		maybeSet = types.NewSetEmpty()
	}

	if maybeSet.Type() != types.ValueTypeSet {
		c.Conn().WriteError(pkg.WrongTypeErr)
		return
	}

	set := maybeSet.Value().(map[string]struct{})

	// We already checked that there are at least 3 arguments.
	// So this should at least iterate once
	result := make([]int, 0)
	for i := 2; i < len(args); i++ {
		_, exists := set[string(args[i])]
		if exists {
			result = append(result, 1)
		} else {
			result = append(result, 0)
		}
	}

	c.Conn().WriteArray(len(result))
	for _, v := range result {
		c.Conn().WriteInt(v)
	}
}
