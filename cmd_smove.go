package redis

import (
	"fmt"
)

// https://redis.io/commands/smove/
// SMOVE source destination member
func SmoveCommand(c *Client, args [][]byte) {
	if len(args) != 4 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()

	sourceKey := string(args[1])
	destinationKey := string(args[2])
	memberKey := string(args[3])

	maybeSource, sourceTtl := db.GetOrExpire(sourceKey, true)

	if maybeSource == nil {
		c.Conn().WriteInt(0)
		return
	} else if maybeSource.Type() != ValueTypeSet {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	sourceSet := maybeSource.(*Set)

	maybeDest, destTtl := db.GetOrExpire(destinationKey, true)

	if maybeDest == nil {
		maybeDest = NewSetEmpty()
	} else if maybeDest.Type() != ValueTypeSet {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	destSet := maybeDest.(*Set)

	existed := sourceSet.RemoveMember(memberKey)

	if !existed {
		c.Conn().WriteInt(0)
		return
	}

	destSet.AddMember(memberKey)

	db.Set(sourceKey, sourceSet, sourceTtl)
	db.Set(destinationKey, destSet, destTtl)

	c.Conn().WriteInt(1)
}
