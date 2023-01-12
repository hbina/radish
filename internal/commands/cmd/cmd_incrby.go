package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/incrby/
func IncrByCommand(c *pkg.Client, args [][]byte) {
	if len(args) != 3 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()

	key := string(args[1])
	incrBy, err := strconv.ParseInt(string(args[2]), 10, 64)

	if err != nil {
		c.Conn().WriteError(util.InvalidIntErr)
		return
	}

	item, exists := db.Storage[key]

	if !exists {
		db.Set(key, types.NewString(fmt.Sprintf("%d", incrBy)), time.Time{})
		c.Conn().WriteInt64(incrBy)
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

	intValue += incrBy

	db.Set(key, types.NewString(fmt.Sprint(intValue)), time.Time{})
	c.Conn().WriteInt64(intValue)
}
