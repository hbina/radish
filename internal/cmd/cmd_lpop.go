package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
)

// https://redis.io/commands/lpop/
func LPopCommand(c *pkg.Client, args [][]byte) {
	if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, args[0]))
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

	l := item.(*List)
	value, valid := l.LPop()

	if valid {
		c.Conn().WriteBulkString(value)
	} else {
		db.Delete(key)
		c.Conn().WriteNull()
	}
}
