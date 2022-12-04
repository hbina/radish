package redis

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// https://redis.io/commands/zunion/
// ZUNION numkeys key [key ...] [WEIGHTS weight [weight ...]] [AGGREGATE <SUM | MIN | MAX>] [WITHSCORES]
func ZunionCommand(c *Client, args [][]byte) {
	implZunionCommand(c, args, false)
}

// Shared generic function for Zunion* family of functions
// https://redis.io/commands/zunion/
// https://redis.io/commands/zunionstore/
func implZunionCommand(c *Client, args [][]byte, store bool) {
	offsetToKeys := 2
	offsetToNumKeys := 1

	// Store requires 1 more argument for thet destination
	if store {
		offsetToKeys += 1
		offsetToNumKeys += 1
	}

	// Check if we have the minimum number of args
	if len(args) < offsetToKeys {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	destination := string(args[1])
	numKeyStr := string(args[offsetToNumKeys])
	numKey64, err := strconv.ParseInt(numKeyStr, 10, 32)

	if err != nil {
		c.Conn().WriteError(InvalidIntErr)
		return
	}

	numKey := int(numKey64)

	// Verify there's enough args for the given numKey
	if len(args)-offsetToKeys < numKey {
		c.Conn().WriteError(SyntaxErr)
		return
	}

	// Collect keys
	keys := make([]string, 0, numKey)
	for i := 0; i < numKey; i++ {
		keys = append(keys, string(args[i+offsetToKeys]))
	}

	// Verify there's enough args for the given numKey
	if len(args)-offsetToKeys < numKey {
		c.Conn().WriteError(SyntaxErr)
		return
	}

	// Collect weights
	weights := make([]float64, 0, numKey)

	// There are still more arguments, lets see if its weight
	// Check if there are enough arguments for weight after subtracting from:
	// 1. the initial offsetIntoNumKeys required arguments (ZUNION destination numkeys)
	// 2. the number of keys provided (numKeys)
	// 3. and the tag for weights (WEIGHTS)
	if len(args)-offsetToKeys-numKey-1 >= numKey {
		if strings.ToLower(string(args[offsetToKeys+numKey])) == "weights" {
			for i := 0; i < numKey; i++ {
				float, err := strconv.ParseFloat(string(args[i+offsetToKeys+numKey+1]), 64)

				if err != nil {
					c.Conn().WriteError(SyntaxErr)
					return
				}

				weights = append(weights, float)
			}

			if len(keys) != len(weights) {
				c.Conn().WriteError(SyntaxErr)
				return
			}
		}
	}

	// What to do with overlapping elements of 2 sets
	// -1 -> not set
	// 0  -> SUM
	// 1  -> MIN
	// 2  -> MAX
	aggregateMode := -1
	withScores := false

	// Parse additional option
	for i := offsetToKeys + len(keys) + 1 + len(weights); i < len(args); i++ {
		switch strings.ToLower(string(args[i])) {
		default:
			{
				c.Conn().WriteError(SyntaxErr)
				return
			}
		case "aggregate":
			{
				// Requires 1 more argument and check if
				// we have found aggregate option before
				if i+1 >= len(args) || aggregateMode != -1 {
					c.Conn().WriteError(SyntaxErr)
					return
				}

				switch strings.ToLower(string(args[i+1])) {
				case "sum":
					{
						aggregateMode = 0
					}
				case "min":
					{
						aggregateMode = 1
					}
				case "max":
					{
						aggregateMode = 2
					}
				}
			}
		case "withscores":
			{
				if store {
					c.Conn().WriteError(SyntaxErr)
					return
				}

				withScores = true
			}
		}
	}

	// If user haven't set aggregate mode yet,
	// then we default to 0
	if aggregateMode == -1 {
		aggregateMode = 0
	}

	db := c.Db()
	union := NewZSet()

	for idx, key := range keys {
		weight := 1.0
		if len(weights) != 0 {
			weight = weights[idx]
		}
		maybeSet, _ := db.GetOrExpire(key, true)

		// If the other set is nil, then the union is no-op
		if maybeSet == nil {
			continue
		} else if maybeSet.Type() != ValueTypeZSet {
			c.Conn().WriteError(WrongTypeErr)
			return
		}

		set := maybeSet.(*ZSet)

		union.inner = *union.inner.Union(&set.inner, aggregateMode, weight)
	}

	if store {
		db.Set(destination, union, time.Time{})
		c.Conn().WriteInt(union.Len())
	} else {
		if withScores {
			c.Conn().WriteArray(union.Len() * 2)
			for key, node := range union.inner.dict {
				c.Conn().WriteBulkString(key)
				c.Conn().WriteBulkString(fmt.Sprint(node))
			}
		} else {
			c.Conn().WriteArray(union.Len())
			for key := range union.inner.dict {
				c.Conn().WriteBulkString(key)
			}
		}
	}
}
