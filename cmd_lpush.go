package redis

import (
	"fmt"

	"github.com/tidwall/redcon"
)

func LPushCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
		return
	} else if len(cmd.Args) == 1 {
		c.Conn().WriteError(fmt.Sprintf("wrong number of arguments for '%s' command", cmd.Args[0]))
		return
	}
	key := string(cmd.Args[1])
	db := c.Db()
	value := db.GetOrExpire(&key, true)
	if value == nil {
		value = NewList()
		db.Set(&key, value, nil)
	} else if value.Type() != ListType {
		c.Conn().WriteError(fmt.Sprintf("%s: key is a %s not a %s", WrongTypeErr, value.TypeFancy(), ListTypeFancy))
		return
	}

	fmt.Println("LPUSH", value)

	list := value.(*List)
	var length int
	c.Redis().Mu().Lock()
	for j := 2; j < len(cmd.Args); j++ {
		v := string(cmd.Args[j])
		length = list.LPush(&v)
	}
	c.Redis().Mu().Unlock()

	c.Conn().WriteInt(length)
}
