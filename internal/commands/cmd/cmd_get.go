package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/get/
func GetCommand(c *pkg.Client, args [][]byte) {
	if len(args) == 1 {
		c.WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	item, _ := c.Db().Get(key)

	if item == nil {
		c.WriteNull()
		return
	}

	if item.Type() == types.ValueTypeString {
		v := item.Value().(string)
		c.WriteBulkString(v)
		return
	} else {
		c.WriteError(fmt.Sprintf("%s: key is a %s not a %s", util.WrongTypeErr, item.TypeFancy(), types.ValueTypeFancyString))
		return
	}
}
