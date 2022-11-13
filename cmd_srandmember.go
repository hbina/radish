package redis

import (
	"fmt"
	"strconv"
)

// https://redis.io/commands/srandmember/
// SRANDMEMBER key [count]
func SrandmemberCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError(ZeroArgumentErr)
		return
	} else if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	var count *int = nil

	if len(args) == 3 {
		count64, err := strconv.ParseInt(string(args[2]), 10, 32)

		if err != nil {
			c.Conn().WriteError(InvalidIntErr)
		}

		count32 := int(count64)
		count = &count32
	}

	db := c.Db()

	maybeSet, _ := db.GetOrExpire(key, true)

	// If any of the sets are nil, then the intersections must be 0
	if maybeSet == nil {
		c.Conn().WriteNull()
		return
	} else if maybeSet.Type() != ValueTypeSet {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	set := maybeSet.(*Set)

	// TODO: If count larger than set, then just delete set
	if count != nil {
		removed := make([]string, 0)

		if *count < 0 {
			for i := 0; i < -*count; i++ {
				member := set.GetRandomMember()

				if member != nil {
					removed = append(removed, *member)
				}
			}
		} else {
			size := 0

			set.ForEachF(func(a string) {
				if size < *count {
					removed = append(removed, a)
				}
				size++
			})
		}

		c.Conn().WriteArray(len(removed))
		for _, k := range removed {
			c.Conn().WriteBulkString(k)
		}
	} else {
		member := set.GetRandomMember()

		if member != nil {
			c.Conn().WriteBulkString(*member)
		} else {
			c.Conn().WriteNull()
		}
	}
}
