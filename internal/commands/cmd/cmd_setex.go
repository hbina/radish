package cmd

import (
	"fmt"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/setex/
// SETEX key seconds value
// This is equivalent to calling `SET key value EX seconds`
func SetexCommand(c *pkg.Client, args [][]byte) {
	if len(args) != 4 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	seconds := string(args[2])
	value := string(args[3])

	newTtl, err := util.ParseTtlFromUnitTime(seconds, int64(time.Second))

	if err != nil {
		c.Conn().WriteError(util.InvalidIntErr)
		return
	}

	db := c.Db()

	db.Set(key, types.NewString(value), newTtl)

	c.Conn().WriteString("OK")
}
