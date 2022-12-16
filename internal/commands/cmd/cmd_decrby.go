package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/decrby/
func DecrByCommand(c *pkg.Client, args [][]byte) {
	if len(args) != 3 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()

	key := string(args[1])
	decrBy, err := strconv.ParseInt(string(args[2]), 10, 64)

	if err != nil {
		c.Conn().WriteError(util.InvalidIntErr)
		return
	}

	item, _ := db.Get(key)

	if item == nil {
		db.Set(key, types.NewString(fmt.Sprintf("%d", decrBy)), time.Time{})
		c.Conn().WriteInt64(decrBy)
		return
	}

	value, ok := item.Value().(string)

	if !ok {
		c.Conn().WriteError(util.WrongTypeErr)
		return
	}

	intValue, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		c.Conn().WriteError(util.InvalidIntErr)
		return
	}

	intValue -= decrBy

	db.Set(key, types.NewString(fmt.Sprint(intValue)), time.Time{})
	c.Conn().WriteInt64(intValue)
}
