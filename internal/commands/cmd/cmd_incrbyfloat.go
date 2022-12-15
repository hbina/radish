package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
)

// https://redis.io/commands/incrbyfloat/
func IncrByFloatCommand(c *pkg.Client, args [][]byte) {
	if len(args) != 3 {
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, args[0]))
		return
	}

	db := c.Db()

	key := string(args[1])
	incrBy, err := strconv.ParseFloat(string(args[2]), 64)

	if err != nil {
		c.Conn().WriteError(pkg.InvalidFloatErr)
		return
	}

	item, exists := db.Storage[key]

	if !exists {
		incrByStr := strconv.FormatFloat(incrBy, 'f', -1, 64)
		db.Set(key, types.NewString(incrByStr), time.Time{})
		c.Conn().WriteString(fmt.Sprintf("\"%s\"", incrByStr))
		return
	}

	value, ok := item.Value().(string)

	if !ok {
		c.Conn().WriteError(pkg.WrongTypeErr)
		return
	}

	floatValue, err := strconv.ParseFloat(value, 64)

	if err != nil {
		c.Conn().WriteError(pkg.InvalidFloatErr)
		return
	}

	floatValue += incrBy

	floatValueStr := strconv.FormatFloat(floatValue, 'f', -1, 64)
	db.Set(key, types.NewString(floatValueStr), time.Time{})
	c.Conn().WriteString(fmt.Sprintf("\"%s\"", floatValueStr))
}
