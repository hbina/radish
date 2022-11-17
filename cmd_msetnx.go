package redis

import (
	"fmt"
	"time"
)

// https://redis.io/commands/msetnx/
// MSETNX key value [key value ...]
func MsetnxCommand(c *Client, args [][]byte) {
	if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	} else if len(args)%2 != 1 {
		// If the number of arguments (excluding the command name) is not even,
		// return syntax error
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, string(args[0])))
		return
	}

	db := c.Db()

	// Returns an error if _any_ of the keys already exist
	for i := 1; i < len(args); i += 2 {
		key := string(args[i])

		if db.Exists(key) {
			c.Conn().WriteInt(0)
			return
		}
	}

	for i := 1; i < len(args); i += 2 {
		keyStr := string(args[i])
		valueStr := string(args[i+1])

		db.Set(keyStr, NewString(valueStr), time.Time{})
	}

	c.Conn().WriteInt(1)
}
