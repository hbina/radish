package cmd

import (
	"fmt"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/setnx/
// SETNX key value
// This is equivalent to calling SET key value NX
func SetNxCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 3 {
		c.WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	value := string(args[2])

	db := c.Db()
	exists := db.Exists(key)

	if exists {
		c.WriteInt(0)
		return
	}

	db.Set(key, types.NewString(value), time.Time{})

	c.WriteInt(1)
}
