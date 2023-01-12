package cmd

import (
	"fmt"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/ttl/
// TTL key
func TtlCommand(c *pkg.Client, args [][]byte) {
	if len(args) != 2 {
		c.WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()
	key := string(args[1])

	item, ttl := db.Get(key)

	if item == nil {
		c.WriteInt(-2)
		return
	} else if item != nil && time.Time.IsZero(ttl) {
		c.WriteInt(-1)
		return
	}

	c.WriteInt64(int64(time.Until(ttl).Seconds()))
}
