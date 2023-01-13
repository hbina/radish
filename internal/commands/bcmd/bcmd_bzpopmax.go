package bcmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/bzmpopmax/
// BZPOPMAX key [key ...] timeout
func BzpopmaxCommand(c *pkg.Client, args [][]byte) *pkg.BlockedCommand {
	if len(args) < 3 {
		c.Conn().WriteError(util.SyntaxErr)
		return nil
	}

	timeoutStr := string(args[len(args)-1])

	timeout64, err := strconv.ParseFloat(timeoutStr, 64)

	if err != nil || timeout64 < 0 {
		c.Conn().WriteError(util.SyntaxErr)
		return nil
	}

	ttl := time.Time{}

	if timeout64 > 0 {
		ttl = time.Now().Add(time.Duration(timeout64 * float64(time.Second)))
	}

	keys := make([]string, 0, len(args)-2)

	for i := 1; i < len(args)-1; i++ {
		keys = append(keys, string(args[i]))
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

		n := set.RemoveByRank(set.Len())

		db.Set(key, types.NewZSetFromSs(set), ttl)

		c.Conn().WriteArray(3)
		c.Conn().WriteBulkString(key)
		c.Conn().WriteBulkString(n.Key)
		c.Conn().WriteBulkString(fmt.Sprint(n.Score))

		return nil
	}

	return pkg.NewBlockedCommand(
		c,
		args,
		ttl,
		time.Duration(timeout64*float64(time.Second)))
}
