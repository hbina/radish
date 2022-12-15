package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
)

// https://redis.io/commands/incr/
func IncrCommand(c *pkg.Client, args [][]byte) {
	if len(args) == 1 {
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()

	key := string(args[1])

	item, exists := db.Storage[key]

	if !exists {
		db.Set(key, types.NewString("1"), time.Time{})
		c.Conn().WriteInt64(1)
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

	intValue++

	db.Set(key, types.NewString(fmt.Sprint(intValue)), time.Time{})
	c.Conn().WriteInt64(intValue)
}
