package redis

import (
	"fmt"
	"go-redis/ref"
	"strconv"

	"github.com/tidwall/redcon"
)

func DecrByFloatCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
		return
	} else if len(cmd.Args) != 3 {
		c.Conn().WriteError(fmt.Sprintf("wrong number of arguments for '%s' command", cmd.Args[0]))
		return
	}

	db := c.Db()

	key := string(cmd.Args[1])
	decrBy, err := strconv.ParseFloat(string(cmd.Args[2]), 64)

	if err != nil {
		c.Conn().WriteError(InvalidFloatErr)
		return
	}

	db.Mu().Lock()
	defer db.Mu().Unlock()

	item, exists := db.keys[key]

	if !exists {
		decrByStr := strconv.FormatFloat(decrBy, 'f', -1, 64)
		db.Set(&key, NewString(ref.String(decrByStr)), nil)
		c.conn.WriteString(fmt.Sprintf("\"%s\"", decrByStr))
		return
	}

	value, ok := item.Value().(*string)

	if !ok {
		c.conn.WriteError(WrongTypeErr)
		return
	}

	floatValue, err := strconv.ParseFloat(*value, 64)

	if err != nil {
		c.conn.WriteError(InvalidFloatErr)
		return
	}

	floatValue -= decrBy

	floatValueStr := strconv.FormatFloat(floatValue, 'f', -1, 64)
	db.Set(&key, NewString(ref.String(floatValueStr)), nil)
	c.conn.WriteString(fmt.Sprintf("\"%s\"", floatValueStr))
}
