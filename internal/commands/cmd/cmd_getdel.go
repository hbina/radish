package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/getdel/
// GETDEL key
func GetdelCommand(c *pkg.Client, args [][]byte) {
	if len(args) != 2 {
		c.Conn().WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	db := c.Db()
	item, _ := db.Get(key)

	if item == nil {
		c.Conn().WriteNull()
		return
	}

	if item.Type() == types.ValueTypeString {
		v := item.Value().(string)
		c.Conn().WriteBulkString(v)
		// Only delete the key if the operation is succesfull
		db.Delete(key)
		return
	} else {
		c.Conn().WriteError(fmt.Sprintf("%s: key is a %s not a %s", util.WrongTypeErr, item.TypeFancy(), types.ValueTypeFancyString))
		return
	}
}
