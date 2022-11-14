package redis

import (
	"fmt"
	"time"
)

// https://redis.io/commands/setx/
// SETX key value
func SetXCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
		return
	} else if len(args) < 3 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	key := string(args[1])
	value := string(args[2])

	genericSetCommand(c, key, value, time.Time{}, SetWriteXx, false)
}
