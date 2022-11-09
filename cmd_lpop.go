package redis

import (
	"fmt"
)

func LPopCommand(c *Client, args [][]byte) {
	if len(args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, "lpop"))
		return
	}
	key := string(args[1])

	db := c.Db()
	i := db.GetOrExpire(&key, true)
	if i == nil {
		c.Conn().WriteNull()
		return
	} else if i.Type() != ValueTypeList {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	l := i.(*List)
	v, b := l.LPop()
	if b {
		db.Delete(&key)
	}

	c.Conn().WriteBulkString(*v)
}
