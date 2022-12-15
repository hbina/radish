package bcmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/bzmpopmin/
// BZPOPMIN key [key ...] timeout
func BzmpopminCommand(c *pkg.Client, args [][]byte) int {
	if len(args) < 4 {
		c.Conn().WriteError(util.SyntaxErr)
		return pkg.BCMD_OK
	}

	timeoutStr := string(args[1])

	timeout64, err := strconv.ParseInt(timeoutStr, 10, 32)

	if err != nil || timeout64 < 0 {
		c.Conn().WriteError(util.SyntaxErr)
		return pkg.BCMD_OK
	}

	fmt.Println(timeout64)

	numKeyStr := string(args[2])

	numKey64, err := strconv.ParseInt(numKeyStr, 10, 32)

	if err != nil || numKey64 < 0 {
		c.Conn().WriteError(util.SyntaxErr)
		return pkg.BCMD_OK
	}

	numKey := int(numKey64)

	if len(args) < 3+numKey {
		c.Conn().WriteError(util.SyntaxErr)
		return pkg.BCMD_OK
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
		c.Conn().WriteError(util.SyntaxErr)
		return pkg.BCMD_OK
	}

	modeStr := strings.ToLower(string(args[3+numKey+1]))

	if modeStr == "min" {
		mode = 0
	} else if modeStr == "max" {
		mode = 1
	} else {
		c.Conn().WriteError(util.SyntaxErr)
		return pkg.BCMD_OK
	}

	// Parse options
	// -1 -> not set
	count := -1

	for i := 3 + numKey + 1; i < len(args); i++ {
		arg := strings.ToLower(string(args[i]))
		switch arg {
		default:
			{
				c.Conn().WriteError(util.SyntaxErr)
				return pkg.BCMD_OK
			}
		case "count":
			{
				if count != -1 {
					c.Conn().WriteError(util.SyntaxErr)
					return pkg.BCMD_OK
				}

				// Need 1 more argument
				if i+1 >= len(args) {
					c.Conn().WriteError(util.SyntaxErr)
					return pkg.BCMD_OK
				}

				i++
				countStr := string(args[i])

				count64, err := strconv.ParseInt(countStr, 10, 32)

				if err != nil {
					c.Conn().WriteError(util.SyntaxErr)
					return pkg.BCMD_OK
				}

				if count64 <= 0 {
					c.Conn().WriteError(util.SyntaxErr)
					return pkg.BCMD_OK
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
		maybeSet, ttl := db.Get(key)

		if maybeSet == nil {
			continue
		}

		if maybeSet.Type() != types.ValueTypeZSet {
			continue
		}

		set := maybeSet.Value().(*types.SortedSet)

		if count > set.Len() {
			count = set.Len()
		}

		var res []*types.SortedSetNode
		options := types.DefaultRangeOptions()

		if mode == 0 {
			res = set.GetRangeByRank(1, count, options)
		} else {
			options.Reverse = true
			res = set.GetRangeByRank(set.Len()+1-count, set.Len(), options)
		}

		for _, n := range res {
			set.Remove(n.Key)
		}

		db.Set(key, types.NewZSetFromSs(set), ttl)

		c.Conn().WriteArray(2)
		c.Conn().WriteBulkString(key)
		c.Conn().WriteArray(len(res))
		for _, n := range res {
			c.Conn().WriteArray(2)
			c.Conn().WriteBulkString(n.Key)
			c.Conn().WriteBulkString(fmt.Sprint(n.Score))
		}

		return pkg.BCMD_OK
	}

	return pkg.BCMD_RETRY
}
