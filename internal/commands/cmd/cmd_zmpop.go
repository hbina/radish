package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/zmpop/
// ZMPOP numkeys key [key ...] <MIN | MAX> [COUNT count]
func ZmpopCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 4 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	numKeyStr := string(args[1])

	numKey64, err := strconv.ParseInt(numKeyStr, 10, 32)

	if err != nil || numKey64 <= 0 {
		c.Conn().WriteError(fmt.Sprintf(util.NegativeIntErr, "numkeys"))
		return
	}

	numKey := int(numKey64)

	if len(args) < 2+numKey {
		c.Conn().WriteError(util.SyntaxErr)
		return
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

	if 2+numKey >= len(args) {
		c.Conn().WriteError(util.SyntaxErr)
		return
	}

	modeStr := strings.ToLower(string(args[2+numKey]))

	if modeStr == "min" {
		mode = 0
	} else if modeStr == "max" {
		mode = 1
	} else {
		c.Conn().WriteError(util.SyntaxErr)
		return
	}

	for i := 2 + numKey + 1; i < len(args); i++ {
		arg := strings.ToLower(string(args[i]))
		switch arg {
		default:
			{
				c.Conn().WriteError(util.SyntaxErr)
				return
			}
		case "count":
			{
				if count != -1 {
					c.Conn().WriteError(util.SyntaxErr)
					return
				}

				// Need 1 more argument
				if i+1 >= len(args) {
					c.Conn().WriteError(util.SyntaxErr)
					return
				}

				i++
				countStr := string(args[i])

				count64, err := strconv.ParseInt(countStr, 10, 32)

				if err != nil || count64 <= 0 {
					c.Conn().WriteError("ERR count must be greater than 0")
					return
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
		maybeSet, ttl := db.Get(key)

		if maybeSet == nil {
			continue
		}

		if maybeSet.Type() != types.ValueTypeZSet {
			c.Conn().WriteError(util.WrongTypeErr)
			return
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

		if c.R3 {
			c.WriteToConn(res, true)
		} else {
			c.Conn().WriteArray(len(res))

			for _, node := range res {
				c.Conn().WriteArray(2)
				c.Conn().WriteBulkString(node.Key)
				c.Conn().WriteBulkString(fmt.Sprint(node.Score))
			}
		}

		return
	}

	if c.R3 {
		c.Conn().WriteNull()
	} else {
		c.Conn().WriteNullArray()
	}
}
