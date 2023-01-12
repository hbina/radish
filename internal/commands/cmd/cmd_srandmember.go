package cmd

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/srandmember/
// SRANDMEMBER key [count]
func SrandmemberCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 2 {
		c.WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	useCount := false
	count := 0

	if len(args) == 3 {
		count64, err := strconv.ParseInt(string(args[2]), 10, 32)

		if err != nil {
			c.WriteError(util.InvalidIntErr)
		}

		useCount = true
		count = int(count64)
	}

	db := c.Db()

	maybeSet, _ := db.Get(key)

	// If any of the sets are nil, then the intersections must be 0
	if maybeSet == nil {
		c.WriteArray(0)
		return
	} else if maybeSet.Type() != types.ValueTypeSet {
		c.WriteError(util.WrongTypeErr)
		return
	}

	set := maybeSet.(*types.Set)

	if useCount {
		if count > set.Len() {
			count = set.Len()
		}

		result := make([]string, 0)

		members := set.GetMembers()
		if count < 0 {
			// If negative, we just append freely to the result.
			// However, we need to check if the set is zero because
			// rand.Intn is inclusive on the LHS.
			for i := 0; i < -count && set.Len() > 0; i++ {
				randomIdx := rand.Intn(len(members))
				member := members[randomIdx]
				result = append(result, member)
			}
		} else {
			set2 := make(map[string]struct{}, count)

			// Try to fairly choose members
			for i := 0; i < (count)*10; i++ {
				randomIdx := rand.Intn(len(members))
				member := members[randomIdx]
				set2[member] = struct{}{}

				if len(set2) == count {
					break
				}
			}

			// If we failed to fill up the result then
			// we just loop over the set.
			// Go says that iterating over map is semi-random
			// This check will always be enough to fill up result because
			// we already truncated count to be at most as big as set
			set.ForEachF(func(k string) bool {
				if len(set2) == count {
					return false
				}
				set2[k] = struct{}{}
				return true
			})

			// Finally, map set to result
			for k := range set2 {
				result = append(result, k)
			}
		}

		c.WriteArray(len(result))
		for _, k := range result {
			c.WriteBulkString(k)
		}
	} else {
		member := set.GetRandomMember()

		if member != nil {
			c.WriteBulkString(*member)
		} else {
			c.WriteNull()
		}
	}
}
