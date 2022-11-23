package redis

import "fmt"

// https://redis.io/commands/get/
func GetCommand(c *Client, args [][]byte) {
	if len(args) == 1 {
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
