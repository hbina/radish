package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/incr/
func IncrCommand(c *pkg.Client, args [][]byte) {
	if len(args) == 1 {
		c.WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()

	key := string(args[1])

	item, exists := db.Storage[key]

	if !exists {
		db.Set(key, types.NewString("1"), time.Time{})
		c.WriteInt64(1)
		return
	}

	value, ok := item.Value().(string)

	if !ok {
		c.WriteError(util.WrongTypeErr)
		return
	}

	intValue, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		c.WriteError(util.InvalidIntErr)
		return
	}

	intValue++

	db.Set(key, types.NewString(fmt.Sprint(intValue)), time.Time{})
	c.WriteInt64(intValue)
}
