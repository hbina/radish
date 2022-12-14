package cmd

import (
	"fmt"
	"strconv"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/lrange/
func LRangeCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])

	start, err := strconv.Atoi(string(args[2]))
	if err != nil {
		c.Conn().WriteError(fmt.Sprintf("%s: %s", util.InvalidIntErr, err.Error()))
		return
	}

	end, err := strconv.Atoi(string(args[3]))
	if err != nil {
		c.Conn().WriteError(fmt.Sprintf("%s: %s", util.InvalidIntErr, err.Error()))
		return
	}

	db := c.Db()
	item, _ := db.Get(key)

	if item == nil {
		c.Conn().WriteArray(0)
		return
	} else if item.Type() != types.ValueTypeList {
		c.Conn().WriteError(util.WrongTypeErr)
		return
	}

	l := item.(*types.List)
	values := l.LRange(start, end)

	c.Conn().WriteArray(len(values))
	for _, v := range values {
		c.Conn().WriteBulkString(v)
	}
}
