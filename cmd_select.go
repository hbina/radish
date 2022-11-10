package redis

import (
	"fmt"
	"strconv"
)

func SelectCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
		return
	} else if len(args) == 1 {
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
