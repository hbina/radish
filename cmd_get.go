package redis

import (
	"fmt"

	"github.com/tidwall/redcon"
)

func GetCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
		return
	} else if len(cmd.Args) == 1 {
		c.Conn().WriteError(fmt.Sprintf("wrong number of arguments for '%s' command", cmd.Args[0]))
		return
	}

	key := string(cmd.Args[1])

	item := c.Db().GetOrExpire(&key, true)
	if item == nil {
		c.Conn().WriteNull()
		return
	}

	if item.Type() != StringType {
		c.Conn().WriteError(fmt.Sprintf("%s: key is a %s not a %s", WrongTypeErr, item.TypeFancy(), StringTypeFancy))
		return
	}

	v := *item.Value().(*string)
	c.Conn().WriteBulkString(v)
}
