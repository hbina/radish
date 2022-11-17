package redis

import (
	"fmt"
	"time"
)

// https://redis.io/commands/mset/
// MSET key value [key value ...]
func MsetCommand(c *Client, args [][]byte) {
	if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	} else if len(args)%2 != 1 {
		// If the number of arguments (excluding the command name) is not even,
		// return syntax error
		c.Conn().WriteError(WrongNumOfArgsErr)
		return
	}

	db := c.Db()

	for i := 1; i < len(args); i += 2 {
		keyStr := string(args[i])
		valueStr := string(args[i+1])

		db.Set(keyStr, NewString(valueStr), time.Time{})
	}

	c.Conn().WriteString("OK")
}
