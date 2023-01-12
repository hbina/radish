package cmd

import (
	"fmt"

	"github.com/hbina/radish/internal/pkg"
	"github.com/hbina/radish/internal/types"
	"github.com/hbina/radish/internal/util"
)

// https://redis.io/commands/lpush/
func LPushCommand(c *pkg.Client, args [][]byte) {
	if len(args) == 1 {
		c.WriteError(fmt.Sprintf(util.WrongNumOfArgsErr, args[0]))
		return
	}
	key := string(args[1])
	db := c.Db()
	value, exp := db.Get(key)

	if value == nil {
		value = types.NewList()
	} else if value.Type() != types.ValueTypeList {
		c.WriteError(fmt.Sprintf("%s: key is a %s not a %s", util.WrongTypeErr, value.TypeFancy(), types.ValueTypeFancyList))
		return
	}

	list := value.(*types.List)
	var length int
	for j := 2; j < len(args); j++ {
		v := string(args[j])
		length = list.LPush(v)
	}
	db.Set(key, list, exp)

	c.WriteInt(length)
}
