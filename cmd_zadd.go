package redis

import (
	"fmt"
	"strconv"
	"time"

	"github.com/zavitax/sortedset-go"
)

// https://redis.io/commands/sadd/
func ZaddCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError(ZeroArgumentErr)
		return
	} else if len(args) < 3 || (len(args)-2)%2 != 0 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])

	// Validate that all the scores are valid floats
	for i := 2; i < len(args); i += 2 {
		_, err := strconv.ParseFloat(string(args[i]), 64)
		if err != nil {
			c.Conn().WriteError(InvalidFloatErr)
			return
		}
	}

	maybeSet := c.Db().Get(key)

	if maybeSet == nil {
		maybeSet = NewZSetEmpty()
	}

	if maybeSet.Type() != ValueTypeZSet {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	set := maybeSet.Value().(sortedset.SortedSet[string, float64, struct{}])

	count := 0
	for i := 2; i < len(args); i += 2 {
		// already validated
		score, _ := strconv.ParseFloat(string(args[i]), 64)
		member := string(args[i+1])
		added := set.AddOrUpdate(member, score, struct{}{})
		if added {
			count++
		}
	}

	c.Db().Set(key, NewZSet(set), time.Time{})

	c.Conn().WriteInt(count)
}
