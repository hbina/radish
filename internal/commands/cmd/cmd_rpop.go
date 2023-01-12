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
		c.WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, "rpop"))
		return
	}

	key := string(args[1])
	db := c.Db()
	item, _ := db.Get(key)

	if item == nil {
		c.WriteNull()
		return
	} else if item.Type() != types.ValueTypeList {
		c.WriteError(util.WrongTypeErr)
		return
	}

	l := item.(*types.List)
	value, valid := l.RPop()

	if valid {
		c.WriteBulkString(value)
	} else {
		db.Delete(key)
		c.WriteNull()
	}
}
