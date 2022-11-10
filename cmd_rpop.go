package redis

import (
	"fmt"
)

func RPopCommand(c *Client, args [][]byte) {
	if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, "rpop"))
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
	value, valid := l.RPop()

	if valid {
		c.Conn().WriteBulkString(value)
	} else {
		db.Delete(key)
		c.Conn().WriteNull()
	}
}
