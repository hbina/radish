package redis

import (
	"fmt"

	"github.com/tidwall/redcon"
)

func LPopCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) < 2 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, "lpop"))
		return
	}
	key := string(cmd.Args[1])

	db := c.Db()
	i := db.GetOrExpire(&key, true)
	if i == nil {
		c.Conn().WriteNull()
		return
	} else if i.Type() != ListType {
		c.Conn().WriteError(WrongTypeErr)
		return
	}

	l := i.(*List)
	c.Redis().Mu().Lock()
	v, b := l.LPop()
	if b {
		db.Delete(&key)
	}
	c.Redis().Mu().Unlock()

	c.Conn().WriteBulkString(*v)
}
