package redis

import (
	"fmt"
	"strconv"

	"github.com/tidwall/redcon"
)

func SelectCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) == 0 {
		c.Conn().WriteError("no argument passed to handler. This should not be possible")
		return
	} else if len(cmd.Args) == 1 {
		c.Conn().WriteError(fmt.Sprintf("wrong number of arguments for '%s' command", cmd.Args[0]))
		return
	}

	index, err := strconv.ParseUint(string(cmd.Args[1]), 10, 32)

	if err != nil {
		c.Conn().WriteError(InvalidIntErr)
	} else {
		c.SelectDb(DatabaseId(index))
		c.Conn().WriteString("OK")
	}
}
