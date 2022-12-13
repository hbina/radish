package cmd

import (
	"fmt"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
)

// https://redis.io/commands/sadd/
func SaddCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 3 {
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

	c.Db().Set(key, NewSetFromMap(set), time.Time{})

	c.Conn().WriteInt(count)
}
