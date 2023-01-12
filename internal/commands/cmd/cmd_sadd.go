package cmd

import (
	"fmt"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/sadd/
func SaddCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 3 {
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

	set := maybeSet.Value().(map[string]struct{})

	// We already checked that there are at least 3 arguments.
	// So this should at least iterate once
	count := 0
	for i := 2; i < len(args); i++ {
		newMember := string(args[i])
		_, found := set[newMember]
		if !found {
			set[newMember] = struct{}{}
			count++
		}

	}

	c.Db().Set(key, types.NewSetFromMap(set), time.Time{})

	c.WriteInt(count)
}
