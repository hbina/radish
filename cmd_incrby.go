package redis

import (
	"fmt"
	"go-redis/ref"
	"strconv"

	"github.com/tidwall/redcon"
)

func IncrByCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
		return
	} else if len(cmd.Args) != 3 {
		c.Conn().WriteError(fmt.Sprintf("wrong number of arguments for '%s' command", cmd.Args[0]))
		return
	}

	db := c.Db()

	key := string(cmd.Args[1])
	incrBy, err := strconv.ParseInt(string(cmd.Args[2]), 10, 64)

	if err != nil {
		c.Conn().WriteError(InvalidIntErr)
		return
	}

	db.Mu().Lock()
	defer db.Mu().Unlock()

	item, exists := db.keys[key]

	if !exists {
		db.Set(&key, NewString(ref.String(fmt.Sprintf("%d", incrBy))), nil)
		c.conn.WriteInt64(incrBy)
		return
	}

	value, ok := item.Value().(*string)

	if !ok {
		c.conn.WriteError(WrongTypeErr)
		return
	}

	intValue, err := strconv.ParseInt(*value, 10, 64)

	if err != nil {
		c.conn.WriteError(InvalidIntErr)
		return
	}

	intValue += incrBy

	db.Set(&key, NewString(ref.String(fmt.Sprint(intValue))), nil)
	c.conn.WriteInt64(intValue)
}
