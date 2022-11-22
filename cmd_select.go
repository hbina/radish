package redis

import (
	"fmt"
	"strconv"
)

// https://redis.io/commands/select/
func SelectCommand(c *Client, args [][]byte) {
	if len(args) == 1 {
		c.Conn().WriteError(fmt.Sprintf(WrongNumOfArgsErr, args[0]))
		return
	}

	index, err := strconv.ParseUint(string(args[1]), 10, 32)

	if err != nil {
		c.Conn().WriteError(InvalidIntErr)
	} else {
		c.SelectDb(DatabaseId(index))
		c.Conn().WriteString("OK")
	}
}
