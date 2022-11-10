package redis

import (
	"fmt"
	"time"
)

func RPushCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
		return
	} else if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	db := c.Db()
	item, _ := db.GetOrExpire(key, true)

	if item == nil {
		item = NewList()
		db.Set(key, item, time.Time{})
	} else if item.Type() != ValueTypeList {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	list := item.(*List)

	var length int
	for j := 2; j < len(args); j++ {
		v := string(args[j])
		length = list.RPush(v)
	}

	c.Conn().WriteInt(length)
}
