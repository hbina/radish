package cmd

import (
	"fmt"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/setx/
// SETX key value
func SetXCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	value := string(args[2])

	db := c.Db()
	exists := db.Exists(key)

	if !exists {
		c.Conn().WriteInt(0)
		return
	}

	db.Set(key, types.NewString(value), time.Time{})

	c.Conn().WriteInt(1)
}
