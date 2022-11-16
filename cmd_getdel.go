package redis

import (
	"fmt"
)

// https://redis.io/commands/getdel/
// GETDEL key
func GetdelCommand(c *Client, args [][]byte) {
	if len(args) != 2 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	db := c.Db()
	item, _ := db.GetOrExpire(key, true)

	if item == nil {
		c.Conn().WriteNull()
		return
	}

	if item.Type() == ValueTypeString {
		v := item.Value().(string)
		c.Conn().WriteBulkString(v)
		// Only delete the key if the operation is succesfull
		db.Delete(key)
		return
	} else {
		c.Conn().WriteError(fmt.Sprintf("%s: key is a %s not a %s", WrongTypeErr, item.TypeFancy(), ValueTypeFancyString))
		return
	}
}
