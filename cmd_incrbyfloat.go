package redis

import (
	"fmt"
	"go-redis/ref"
	"strconv"

	"github.com/tidwall/redcon"
)

func IncrByFloatCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
		return
	} else if len(cmd.Args) != 3 {
		c.Conn().WriteError(fmt.Sprintf("wrong number of arguments for '%s' command", cmd.Args[0]))
		return
	}

	db := c.Db()

	key := string(cmd.Args[1])
	incrBy, err := strconv.ParseFloat(string(cmd.Args[2]), 64)

	if err != nil {
		c.Conn().WriteError(InvalidFloatErr)
		return
	}

	db.Mu().Lock()
	defer db.Mu().Unlock()

	item, exists := db.keys[key]

	if !exists {
		incrByStr := strconv.FormatFloat(incrBy, 'f', -1, 64)
		db.Set(&key, NewString(ref.String(incrByStr)), nil)
		c.conn.WriteString(fmt.Sprintf("\"%s\"", incrByStr))
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

	floatValue += incrBy

	floatValueStr := strconv.FormatFloat(floatValue, 'f', -1, 64)
	db.Set(&key, NewString(ref.String(floatValueStr)), nil)
	c.conn.WriteString(fmt.Sprintf("\"%s\"", floatValueStr))
}
