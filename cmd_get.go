package redis

import (
	"fmt"
)

func GetCommand(c *Client, args [][]byte) {
	GetCommandRaw(c, args)
}

func GetCommandRaw(c *Client, args [][]byte) bool {
	if len(args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
		return false
	} else if len(args) == 1 {
		c.Conn().WriteError(fmt.Sprintf("wrong number of arguments for '%s' command", args[0]))
		return false
	}

	key := string(args[1])

	item := c.Db().GetOrExpire(key, true)
	if item == nil {
		c.Conn().WriteNull()
		return true
	}

	if item.Type() == ValueTypeString {
		v := *item.Value().(*string)
		c.Conn().WriteBulkString(v)
		return true
	} else {
		c.Conn().WriteError(fmt.Sprintf("%s: key is a %s not a %s", WrongTypeErr, item.TypeFancy(), ValueTypeFancyString))
		return false
	}
}
