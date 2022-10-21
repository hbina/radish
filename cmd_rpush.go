package redis

import (
	"fmt"

	"github.com/tidwall/redcon"
)

func RPushCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
		return
	} else if len(cmd.Args) < 3 {
		c.Conn().WriteError(fmt.Sprintf("wrong number of arguments for '%s' command", cmd.Args[0]))
		return
	}

	key := string(cmd.Args[1])
	db := c.Db()
	item := db.GetOrExpire(&key, true)

	if item == nil {
		item = NewList()
		db.Set(&key, item, nil)
	} else if item.Type() != ListType {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	list := item.(*List)

	var length int
	c.Redis().Mu().Lock()
	for j := 2; j < len(cmd.Args); j++ {
		v := string(cmd.Args[j])
		length = list.RPush(&v)
	}
	c.Redis().Mu().Unlock()

	c.Conn().WriteInt(length)
}
