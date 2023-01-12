package cmd

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/sintercard/
// SINTERCARD numkeys key [key ...] [LIMIT limit]
// TODO: Cleanup this mess. It feels like this shouldn't be as complicated as this?
func SintercardCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 3 {
		c.WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	numberOfKeys64, err := strconv.ParseInt(string(args[1]), 10, 32)

	if err != nil {
		c.WriteError(fmt.Sprintf(util.NegativeIntErr, "numkeys"))
		return
	}

	numberOfKeys := int(numberOfKeys64)

	if numberOfKeys <= 0 {
		c.WriteError(fmt.Sprintf(util.NegativeIntErr, "numkeys"))
		return
	}

	// Should not be possible to have more keys than the args passed
	if numberOfKeys > len(args)-2 {
		c.WriteError("ERR Number of keys can't be greater than number of args")
		return
	}

	// The only additional args that can be passed is LIMIT <limit>
	if numberOfKeys != len(args)-2 && numberOfKeys != len(args)-4 {
		c.WriteError(util.SyntaxErr)
		return
	}

	// Collect keys
	// TODO: Can optimize by removing this temporary array and use the args directly.
	keys := make([]string, 0, numberOfKeys)
	for i := 0; i < numberOfKeys; i++ {
		keys = append(keys, string(args[i+2]))
	}

	// Parse limit option
	limit := math.MaxInt

	// number of keys should be equal to the length of the args minus
	// command name, number of keys, limit and limit number
	if len(args)-4 == numberOfKeys {
		limitOption := string(args[len(args)-2])
		limitValue64, err := strconv.ParseInt(string(args[len(args)-1]), 10, 32)

		// TODO: I think this should be a syntax error if its not limit
		if strings.ToLower(limitOption) != "limit" || err != nil || limitValue64 < 0 {
			c.WriteError("ERR LIMIT can't be negative")
			return
		}

		limit = int(limitValue64)
	}

	db := c.Db()

	var intersection *types.Set = nil

	// TODO: Is it possible to optimize using the fact that we know what the
	// upper bound is?
	// TODO: Should be able to optimize this further by breaking early from this loop
	for _, key := range keys {
		maybeSet, _ := db.Get(key)

		// If any of the sets are nil, then the intersections must be 0
		if maybeSet == nil {
			maybeSet = types.NewSetEmpty()
		} else if maybeSet.Type() != types.ValueTypeSet {
			c.WriteError(util.WrongTypeErr)
			return
		}

		set := maybeSet.(*types.Set)

		if intersection == nil {
			intersection = set
		} else {
			intersection = intersection.Intersect(set)
		}
	}

	if intersection == nil {
		c.WriteInt(0)
		return
	}

	if limit > intersection.Len() || limit == 0 {
		c.WriteInt(intersection.Len())
	} else {
		c.WriteInt(limit)
	}
}
