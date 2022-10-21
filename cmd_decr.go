package redis

import (
	"fmt"
	"go-redis/ref"
	"strconv"

	"github.com/tidwall/redcon"
)

func DecrCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
		return
	} else if len(cmd.Args) == 1 {
		c.Conn().WriteError(fmt.Sprintf("wrong number of arguments for '%s' command", cmd.Args[0]))
		return
	}

	db := c.Db()

	key := string(cmd.Args[1])

	db.Mu().Lock()
	defer db.Mu().Unlock()

	item, exists := db.keys[key]

	if !exists {
		db.Set(&key, NewString(ref.String("-1")), nil)
		c.conn.WriteInt64(-1)
		return
	}

	value, ok := item.Value().(*string)

	if !ok {
		c.conn.WriteError(InvalidIntErr)
		return
	}

	intValue, err := strconv.ParseInt(*value, 10, 64)

	if err != nil {
		c.conn.WriteError(WrongTypeErr)
		return
	}

	intValue--

	db.Set(&key, NewString(ref.String(fmt.Sprint(intValue))), nil)
	c.conn.WriteInt64(intValue)
}
