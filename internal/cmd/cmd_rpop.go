package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
)

// https://redis.io/commands/rpop/
func RPopCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, "rpop"))
		return
	}

	key := string(args[1])
	db := c.Db()
	item, _ := db.GetOrExpire(key, true)

	if item == nil {
		c.Conn().WriteNull()
		return
	} else if item.Type() != types.ValueTypeList {
		c.Conn().WriteError(pkg.WrongTypeErr)
		return
	}

	l := item.(*types.List)
	value, valid := l.RPop()

	if valid {
		c.Conn().WriteBulkString(value)
	} else {
		db.Delete(key)
		c.Conn().WriteNull()
	}
}
