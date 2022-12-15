package cmd

import (
	"fmt"
	"strconv"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
)

// https://redis.io/commands/lrange/
func LRangeCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])

	start, err := strconv.Atoi(string(args[2]))
	if err != nil {
		c.Conn().WriteError(fmt.Sprintf("%s: %s", pkg.InvalidIntErr, err.Error()))
		return
	}

	end, err := strconv.Atoi(string(args[3]))
	if err != nil {
		c.Conn().WriteError(fmt.Sprintf("%s: %s", pkg.InvalidIntErr, err.Error()))
		return
	}

	db := c.Db()
	item, _ := db.GetOrExpire(key, true)

	if item == nil {
		c.Conn().WriteArray(0)
		return
	} else if item.Type() != types.ValueTypeList {
		c.Conn().WriteError(pkg.WrongTypeErr)
		return
	}

	l := item.(*types.List)
	values := l.LRange(start, end)

	c.Conn().WriteArray(len(values))
	for _, v := range values {
		c.Conn().WriteBulkString(v)
	}
}
