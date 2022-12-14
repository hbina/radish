package cmd

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
)

// https://redis.io/commands/zunion/
// ZUNION numkeys key [key ...] [WEIGHTS weight [weight ...]] [AGGREGATE <SUM | MIN | MAX>] [WITHSCORES]
func ZunionCommand(c *pkg.Client, args [][]byte) {
	implZSetSetOperationCommand(c, args, false, ZSetOperationUnion, false)
}

const (
	ZSetOperationUnion = iota
	ZSetOperationInter
	ZSetOperationDiff
)

// Shared functions for operations on ZSet
// https://redis.io/commands/zunion/
// https://redis.io/commands/zunionstore/
// https://redis.io/commands/zunioncard/
// https://redis.io/commands/zinter/
// https://redis.io/commands/zinterstore/
// https://redis.io/commands/zintercard/
// https://redis.io/commands/zdiff/
// https://redis.io/commands/zdiffstore/
// https://redis.io/commands/zdiffcard/
func implZSetSetOperationCommand(c *pkg.Client, args [][]byte,
	store bool, operation int, card bool) {
	offsetToKeys := 2
	currArg := 0

	// Store requires 1 more argument for thet destination
	if store {
		offsetToKeys += 1
	}

	// Check if we have the minimum number of args
	if len(args) < offsetToKeys {
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, args[0]))
		return
	}

	destination := string(args[1])
	numKeyStr := string(args[offsetToKeys-1])
	numKey64, err := strconv.ParseInt(numKeyStr, 10, 32)

	if err != nil {
		c.Conn().WriteError(pkg.InvalidIntErr)
		return
	}

	numKey := int(numKey64)

	if numKey < 0 {
		c.Conn().WriteError("ERR LIMIT can't be negative")
		return
	}

	// Verify there's enough args for the given numKey
	if len(args)-offsetToKeys < numKey {
		c.Conn().WriteError(pkg.SyntaxErr)
		return
	}

	// Collect keys
	keys := make([]string, 0, numKey)
	for i := 0; i < numKey; i++ {
		keys = append(keys, string(args[i+offsetToKeys]))
	}
	currArg += offsetToKeys + len(keys)

	// Verify there's enough args for the given numKey
	if len(args)-offsetToKeys < numKey {
		c.Conn().WriteError(pkg.SyntaxErr)
		return
	}

	// Collect weights
	weights := make([]float64, 0, numKey)

	// There are still more arguments, lets see if its weight
	// Check if there are enough arguments for weight after subtracting from:
	// 1. the initial offsetIntoNumKeys required arguments (ZUNION destination numkeys)
	// 2. the number of keys provided (numKeys)
	// 3. and the tag for weights (WEIGHTS)
	// card doesn't care about any of this so we skip for it
	if len(args)-offsetToKeys-numKey-1 >= numKey && !card {
		if strings.ToLower(string(args[offsetToKeys+numKey])) == "weights" {
			for i := 0; i < numKey; i++ {
				float, err := strconv.ParseFloat(string(args[i+offsetToKeys+numKey+1]), 64)

				if err != nil || math.IsNaN(float) {
					c.Conn().WriteError("ERR weight value is not a float")
					return
				}

				weights = append(weights, float)
			}

			if len(keys) != len(weights) {
				c.Conn().WriteError(pkg.SyntaxErr)
				return
			}
		}
	}

	if len(weights) != 0 {
		currArg += 1 + len(weights)
	}

	// What to do with overlapping elements of 2 sets
	// -1 -> not set
	// 0  -> SUM
	// 1  -> MIN
	// 2  -> MAX
	aggregateMode := -1
	withScores := false
	limit := math.MaxInt

	// Parse additional option
	for i := currArg; i < len(args); i++ {
		switch strings.ToLower(string(args[i])) {
		default:
			{
				c.Conn().WriteError(pkg.SyntaxErr)
				return
			}
		case "aggregate":
			{
				if card || operation == ZSetOperationDiff {
					c.Conn().WriteError(pkg.SyntaxErr)
					return
				}

				// Requires 1 more argument and check if
				// we have found aggregate option before
				if i+1 >= len(args) || aggregateMode != -1 {
					c.Conn().WriteError(pkg.SyntaxErr)
					return
				}

				i++
				switch strings.ToLower(string(args[i])) {
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
				if store || card {
					c.Conn().WriteError(pkg.SyntaxErr)
					return
				}

				withScores = true
			}
		case "limit":
			{
				if !card || operation == ZSetOperationDiff {
					c.Conn().WriteError(pkg.SyntaxErr)
					return
				}

				if i+1 >= len(args) {
					c.Conn().WriteError(pkg.SyntaxErr)
					return
				}

				i++
				limitStr := string(args[i])

				limit64, err := strconv.ParseInt(limitStr, 10, 32)

				if err != nil || limit64 < 0 {
					c.Conn().WriteError("ERR LIMIT must be a positive integer")
					return
				}

				if limit64 == 0 {
					limit = math.MaxInt
				} else {
					limit = int(limit64)
				}
			}
		}
	}

	// If user haven't set aggregate mode yet,
	// then we default to 0
	if aggregateMode == -1 {
		aggregateMode = 0
	}

	db := c.Db()
	var result *types.ZSet = nil

	for idx, key := range keys {
		weight := 1.0

		if len(weights) != 0 {
			weight = weights[idx]
		}

		maybeSet, _ := db.GetOrExpire(key, true)

		if maybeSet == nil {
			maybeSet = types.NewZSet()
		}

		var set *types.ZSet = nil

		if maybeSet.Type() == types.ValueTypeZSet {
			set = maybeSet.(*types.ZSet)
		} else if maybeSet.Type() == types.ValueTypeSet {
			set = maybeSet.(*types.Set).ToZSet()
		} else {
			c.Conn().WriteError(pkg.WrongTypeErr)
			return
		}

		if result == nil {
			// We initialize a new set using the weights
			result = types.NewZSet().Union(set, aggregateMode, weight)
		} else {
			if operation == ZSetOperationUnion {
				result = result.Union(set, aggregateMode, weight)
			} else if operation == ZSetOperationInter {
				result = result.Intersect(set, aggregateMode, weight)
			} else if operation == ZSetOperationDiff {
				result = result.Diff(set)
			} else {
				c.Conn().WriteError(pkg.SyntaxErr)
				return
			}
		}

		if card && result != nil && result.Len() >= limit {
			break
		}
	}

	if result == nil {
		result = types.NewZSet()
	}

	if card {
		if result.Len() > limit {
			c.Conn().WriteInt(limit)
		} else {
			c.Conn().WriteInt(result.Len())
		}
	} else if store {
		db.Set(destination, result, time.Time{})
		c.Conn().WriteInt(result.Len())
	} else {
		if withScores {
			c.Conn().WriteArray(result.Len() * 2)
			for _, node := range result.Inner.GetRangeByRank(1, result.Len(), types.DefaultRangeOptions()) {
				c.Conn().WriteBulkString(node.Key)
				c.Conn().WriteBulkString(fmt.Sprint(node.Score))
			}
		} else {
			c.Conn().WriteArray(result.Len())
			for _, node := range result.Inner.GetRangeByRank(1, result.Len(), types.DefaultRangeOptions()) {
				c.Conn().WriteBulkString(node.Key)
			}
		}
	}
}
