package redis

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zavitax/sortedset-go"
)

const (
	ZaddCompareMode = iota
	ZaddCompareGt
	ZaddCompareLt
)

const (
	ZaddExpireMode = iota
	ZaddExpireNx
	ZaddExpireXx
)

// https://redis.io/commands/sadd/
// ZADD key [NX | XX] [GT | LT] [CH] [INCR] score member [score member ...]
func ZaddCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError(ZeroArgumentErr)
		return
	} else if len(args) < 3 || (len(args)-2)%2 != 0 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])

	// Parse options
	optionCount := 0
	compareMode := ZaddCompareMode
	expireMode := ZaddExpireMode
	chEnabled := false
	incrEnabled := false

	// TODO: Can be optimized to end when we encounter an integer
	for i := 2; i < len(args); i++ {
		arg := strings.ToLower(string(args[i]))
		switch arg {
		case "xx":
			{
				if expireMode != ZaddExpireMode {
					c.Conn().WriteError(SyntaxErr)
					return
				}
				expireMode = ZaddExpireXx
				optionCount++
			}
		case "nx":
			{
				if expireMode != ZaddExpireMode {
					c.Conn().WriteError(SyntaxErr)
					return
				}
				expireMode = ZaddExpireNx
				optionCount++
			}
		case "gt":
			{
				if compareMode != ZaddCompareMode {
					c.Conn().WriteError(SyntaxErr)
					return
				}
				compareMode = ZaddCompareGt
				optionCount++
			}
		case "lt":
			{
				if compareMode != ZaddCompareMode {
					c.Conn().WriteError(SyntaxErr)
					return
				}
				compareMode = ZaddCompareLt
				optionCount++
			}
		case "ch":
			{
				chEnabled = true
				optionCount++
			}
		case "incr":
			{
				incrEnabled = true
				optionCount++
			}
		}
	}

	// Validate that all the scores are valid floats
	for i := 2 + optionCount; i < len(args); i += 2 {
		_, err := strconv.ParseFloat(string(args[i]), 64)
		if err != nil {
			c.Conn().WriteError(InvalidFloatErr)
			return
		}
	}
	// Redis does not support multiple score-element pair when doing INCR option
	// for some reasons...
	if incrEnabled && len(args)-optionCount-2 > 2 {
		c.Conn().WriteError(fmt.Sprintf(MultipleElementIncrementPairErr, "INCR"))
		return
	}

	maybeSet := c.Db().Get(key)

	if maybeSet == nil {
		maybeSet = NewZSet()
	}

	if maybeSet.Type() != ValueTypeZSet {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	set := maybeSet.Value().(sortedset.SortedSet[string, float64, struct{}])

	addedCount := 0
	for i := 2; i < len(args); i += 2 {
		// We already validated all scores to be valid
		score, _ := strconv.ParseFloat(string(args[i]), 64)
		member := string(args[i+1])
		old := set.GetByKey(member)

		if (old == nil && expireMode == ZaddExpireNx) ||
			(old != nil && expireMode == ZaddExpireXx) ||
			(old != nil && compareMode == ZaddCompareGt && score <= old.Score()) ||
			(old != nil && compareMode == ZaddCompareLt && score >= old.Score()) {
			continue
		}

		added := set.AddOrUpdate(member, score, struct{}{})
		if added || chEnabled {
			addedCount++
		}
	}

	c.Db().Set(key, NewZSetFromSortedSet(set), time.Time{})

	c.Conn().WriteInt(addedCount)
}
