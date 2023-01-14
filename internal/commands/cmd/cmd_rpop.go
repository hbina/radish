package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/rpop/
func RPopCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, "rpop"))
		return
	}

	key := string(args[1])
	db := c.Db()
	item, _ := db.Get(key)

	if item == nil {
		if c.R3 {
			c.Conn().WriteNull()
		} else {
			c.Conn().WriteNullBulk()
		}
		return
	} else if item.Type() != types.ValueTypeList {
		c.Conn().WriteError(util.WrongTypeErr)
		return
	}

	l := item.(*types.List)
	value, valid := l.RPop()

	if valid {
		c.Conn().WriteBulkString(value)
	} else {
		db.Delete(key)
		if c.R3 {
			c.Conn().WriteNull()
		} else {
			c.Conn().WriteNullBulk()
		}
	}
}
