package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/smove/
// SMOVE source destination member
func SmoveCommand(c *pkg.Client, args [][]byte) {
	if len(args) != 4 {
		c.WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()

	sourceKey := string(args[1])
	destinationKey := string(args[2])
	memberKey := string(args[3])

	maybeSource, sourceTtl := db.Get(sourceKey)

	if maybeSource == nil {
		c.WriteInt(0)
		return
	} else if maybeSource.Type() != types.ValueTypeSet {
		c.WriteError(util.WrongTypeErr)
		return
	}

	sourceSet := maybeSource.(*types.Set)

	maybeDest, destTtl := db.Get(destinationKey)

	if maybeDest == nil {
		maybeDest = types.NewSetEmpty()
	} else if maybeDest.Type() != types.ValueTypeSet {
		c.WriteError(util.WrongTypeErr)
		return
	}

	destSet := maybeDest.(*types.Set)

	existed := sourceSet.RemoveMember(memberKey)

	if !existed {
		c.WriteInt(0)
		return
	}

	destSet.AddMember(memberKey)

	db.Set(sourceKey, sourceSet, sourceTtl)
	db.Set(destinationKey, destSet, destTtl)

	c.WriteInt(1)
}
