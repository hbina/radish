package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/smembers/
func SmembersCommand(c *pkg.Client, args [][]byte) {
	if len(args) != 2 {
		c.WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	maybeSet, _ := c.Db().Get(key)

	if maybeSet == nil {
		maybeSet = types.NewSetEmpty()
	}

	if maybeSet.Type() != types.ValueTypeSet {
		c.WriteError(util.WrongTypeErr)
		return
	}

	set := maybeSet.(*types.Set)

	result := make([]string, 0)
	set.ForEachF(func(k string) bool {
		result = append(result, k)
		return true
	})

	c.WriteArray(len(result))
	for _, v := range result {
		c.WriteBulkString(v)
	}
}
