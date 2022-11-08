package redis

import (
	"fmt"
)

func LPushCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
		return
	} else if len(args) == 1 {
		c.Conn().WriteError(fmt.Sprintf("wrong number of arguments for '%s' command", args[0]))
		return
	}
	key := string(args[1])
	db := c.Db()
	value := db.GetOrExpire(&key, true)
	if value == nil {
		value = NewList()
		db.Set(&key, value, nil)
	} else if value.Type() != ListType {
		c.Conn().WriteError(fmt.Sprintf("%s: key is a %s not a %s", WrongTypeErr, value.TypeFancy(), ListTypeFancy))
		return
	}

	list := value.(*List)
	var length int
	for j := 2; j < len(args); j++ {
		v := string(args[j])
		length = list.LPush(&v)
	}

	c.Conn().WriteInt(length)
}
