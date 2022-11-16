package redis

import "fmt"

func GetCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
		return
	} else if len(args) == 1 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	item, _ := c.Db().GetOrExpire(key, true)

	if item == nil {
		c.Conn().WriteNull()
		return
	}

	if item.Type() == ValueTypeString {
		v := item.Value().(string)
		c.Conn().WriteBulkString(v)
		return
	} else {
		c.Conn().WriteError(fmt.Sprintf("%s: key is a %s not a %s", WrongTypeErr, item.TypeFancy(), ValueTypeFancyString))
		return
	}
}
