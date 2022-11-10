package redis

import (
	"fmt"
)

func LPopCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
		return
	} else if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	db := c.Db()
	item, _ := db.GetOrExpire(key, true)

	if item == nil {
		c.Conn().WriteNull()
		return
	} else if item.Type() != ValueTypeList {
		c.Conn().WriteError(WrongTypeErr)
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
