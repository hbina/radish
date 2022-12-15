package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
)

// https://redis.io/commands/decrby/
func DecrByCommand(c *pkg.Client, args [][]byte) {
	if len(args) != 3 {
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()

	key := string(args[1])
	decrBy, err := strconv.ParseInt(string(args[2]), 10, 64)

	if err != nil {
		c.Conn().WriteError(pkg.InvalidIntErr)
		return
	}

	item, exists := db.Storage[key]

	if !exists {
		db.Set(key, types.NewString(fmt.Sprintf("%d", decrBy)), time.Time{})
		c.Conn().WriteInt64(decrBy)
		return
	}

	value, ok := item.Value().(string)

	if !ok {
		c.Conn().WriteError(pkg.WrongTypeErr)
		return
	}

	intValue, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		c.Conn().WriteError(pkg.InvalidIntErr)
		return
	}

	intValue -= decrBy

	db.Set(key, types.NewString(fmt.Sprint(intValue)), time.Time{})
	c.Conn().WriteInt64(intValue)
}
