package redis

import (
	"fmt"
	"strconv"
	"strings"
)

// https://redis.io/commands/bzmpop/
// BZMPOP timeout numkeys key [key ...] <MIN | MAX> [COUNT count]
// This command should behave exactly like ZMPOP except that it
// will block until it pops a set.
func BzmpopCommand(c *Client, args [][]byte) int {
	if len(args) < 4 {
		c.Conn().WriteError(SyntaxErr)
		return BCMD_OK
	}

	timeoutStr := string(args[1])

	timeout64, err := strconv.ParseInt(timeoutStr, 10, 32)

	if err != nil || timeout64 < 0 {
		c.Conn().WriteError(SyntaxErr)
		return BCMD_OK
	}

	fmt.Println(timeout64)

	numKeyStr := string(args[2])

	numKey64, err := strconv.ParseInt(numKeyStr, 10, 32)

	if err != nil || numKey64 < 0 {
		c.Conn().WriteError(SyntaxErr)
		return BCMD_OK
	}

	numKey := int(numKey64)

	if len(args) < 3+numKey {
		c.Conn().WriteError(SyntaxErr)
		return BCMD_OK
	}

	keys := make([]string, 0, numKey)

	for i := 3; i < 3+numKey; i++ {
		keys = append(keys, string(args[i]))
	}

	// -1 -> not set
	// 0  -> min
	// 1  -> max
	mode := -1

	if len(args) < 3+numKey+1 {
		c.Conn().WriteError(SyntaxErr)
		return BCMD_OK
	}

	modeStr := strings.ToLower(string(args[3+numKey+1]))

	if modeStr == "min" {
		mode = 0
	} else if modeStr == "max" {
		mode = 1
	} else {
		c.Conn().WriteError(SyntaxErr)
		return BCMD_OK
	}

	// Parse options
	// -1 -> not set
	count := -1

	for i := 3 + numKey + 1; i < len(args); i++ {
		arg := strings.ToLower(string(args[i]))
		switch arg {
		default:
			{
				c.Conn().WriteError(SyntaxErr)
				return BCMD_OK
			}
		case "count":
			{
				if count != -1 {
					c.Conn().WriteError(SyntaxErr)
					return BCMD_OK
				}

				// Need 1 more argument
				if i+1 >= len(args) {
					c.Conn().WriteError(SyntaxErr)
					return BCMD_OK
				}

				i++
				countStr := string(args[i])

				count64, err := strconv.ParseInt(countStr, 10, 32)

				if err != nil {
					c.Conn().WriteError(SyntaxErr)
					return BCMD_OK
				}

				if count64 <= 0 {
					c.Conn().WriteError(SyntaxErr)
					return BCMD_OK
				}

				count = int(count64)
			}
		}
	}

	// If not set then default
	if count == -1 {
		count = 1
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

		return BCMD_OK
	}

	return BCMD_RETRY
}
