package redis

import (
	"fmt"
)

// https://redis.io/commands/lpush/
func LPushCommand(c *Client, args [][]byte) {
	if len(args) == 1 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}
	key := string(args[1])
	db := c.Db()
	value, exp := db.GetOrExpire(key, true)

	if value == nil {
		value = NewList()
	} else if value.Type() != ValueTypeList {
		c.Conn().WriteError(fmt.Sprintf("%s: key is a %s not a %s", WrongTypeErr, value.TypeFancy(), ValueTypeFancyList))
		return
	}

	list := value.(*List)
	var length int
	for j := 2; j < len(args); j++ {
		v := string(args[j])
		length = list.LPush(v)
	}
	db.Set(key, list, exp)

	c.Conn().WriteInt(length)
}
