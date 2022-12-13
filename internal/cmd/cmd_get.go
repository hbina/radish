package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
)

// https://redis.io/commands/get/
func GetCommand(c *pkg.Client, args [][]byte) {
	if len(args) == 1 {
		c.Conn().WriteError(fmt.Sprintf(pkg.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	item, _ := c.Db().GetOrExpire(key, true)

	if item == nil {
		c.Conn().WriteNull()
		return
	}

	if item.Type() == types.ValueTypeString {
		v := item.Value().(string)
		c.Conn().WriteBulkString(v)
		return
	} else {
		c.Conn().WriteError(fmt.Sprintf("%s: key is a %s not a %s", pkg.WrongTypeErr, item.TypeFancy(), types.ValueTypeFancyString))
		return
	}
}
