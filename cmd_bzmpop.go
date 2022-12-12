package redis

import (
	"fmt"
	"strconv"
	"strings"
)

// https://redis.io/commands/bzmpop/
// BZMPOP numkeys key [key ...] <MIN | MAX> [COUNT count]
// This command should behave exactly like ZMPOP except that it
// will block until it pops a set.
func BzmpopCommand(c *Client, args [][]byte) int {
	if len(args) < 3 {
		return BCMD_BADARG
	}

	numKeyStr := string(args[1])

	numKey64, err := strconv.ParseInt(numKeyStr, 10, 32)

	if err != nil || numKey64 < 0 {
		return BCMD_BADARG
	}

	numKey := int(numKey64)

	if len(args) < 2+numKey {
		return BCMD_BADARG
	}

	keys := make([]string, 0, numKey)

	for i := 2; i < 2+numKey; i++ {
		keys = append(keys, string(args[i]))
	}

	// Parse options
	// -1 -> not set
	count := -1
	// -1 -> not set
	// 0  -> min
	// 1  -> max
	mode := -1

	for i := 2 + numKey; i < len(args); i++ {
		arg := strings.ToLower(string(args[i]))
		switch arg {
		default:
			{
				return BCMD_BADARG
			}
		case "min":
			{
				if mode != -1 {
					return BCMD_BADARG
				}

				mode = 0
			}
		case "max":
			{
				if mode != -1 {
					return BCMD_BADARG
				}

				mode = 1
			}
		case "count":
			{
				if count != -1 {
					return BCMD_BADARG
				}

				// Need 1 more argument
				if i+1 >= len(args) {
					return BCMD_BADARG
				}

				i++
				countStr := string(args[i])

				count64, err := strconv.ParseInt(countStr, 10, 32)

				if err != nil {
					return BCMD_BADARG
				}

				if count64 <= 0 {
					return BCMD_BADARG
				}

				count = int(count64)
			}
		}
	}

	// If not set then default
	if count == -1 {
		count = 1
	}

	// If not set then default
	if mode == -1 {
		mode = 0
	}

	db := c.Db()

	for _, key := range keys {
		maybeSet, ttl := db.GetOrExpire(key, true)

		if maybeSet == nil {
			continue
		}

		if maybeSet.Type() != ValueTypeZSet {
			continue
		}

		set := maybeSet.Value().(*SortedSet)

		if count > set.Len() {
			count = set.Len()
		}

		var res []*SortedSetNode
		options := DefaultRangeOptions()

		if mode == 0 {
			res = set.GetRangeByRank(1, count, options)
		} else {
			options.reverse = true
			res = set.GetRangeByRank(set.Len()+1-count, set.Len(), options)
		}

		for _, n := range res {
			set.Remove(n.key)
		}

		db.Set(key, NewZSetFromSs(set), ttl)

		c.Conn().WriteArray(2)
		c.Conn().WriteBulkString(key)
		c.Conn().WriteArray(len(res))
		for _, n := range res {
			c.Conn().WriteArray(2)
			c.Conn().WriteBulkString(n.key)
			c.Conn().WriteBulkString(fmt.Sprint(n.score))
		}
	}

	return BCMD_RETRY
}
