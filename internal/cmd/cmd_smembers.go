package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
)

// https://redis.io/commands/smembers/
func SmembersCommand(c *pkg.Client, args [][]byte) {
	if len(args) != 2 {
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	maybeSet := c.Db().Get(key)

	if maybeSet == nil {
		maybeSet = NewSetEmpty()
	}

	if maybeSet.Type() != types.ValueTypeSet {
		c.Conn().WriteError(pkg.WrongTypeErr)
		return
	}

	set := maybeSet.(*types.Set)

	result := make([]string, 0)
	set.ForEachF(func(k string) {
		result = append(result, k)
	})

	c.Conn().WriteArray(len(result))
	for _, v := range result {
		c.Conn().WriteBulkString(v)
	}
}
