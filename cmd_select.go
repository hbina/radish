package redis

import (
	"strconv"

	"github.com/tidwall/redcon"
)

func SelectCommand(c *Client, cmd redcon.Command) {
	if len(cmd.Args) < 2 {
		c.Conn().WriteError("wrong number of arguments for 'select' command")
	}
	index, err := strconv.ParseUint(string(cmd.Args[1]), 10, 32)
	if err != nil {
		c.Conn().WriteError("value is not an integer or out of range")
	} else {
		// SAFETY: Conversion here is safe because we require the conversion to fit
		// in 32-bit unsigned integer.
		c.SelectDb(DatabaseId(index))
		c.Conn().WriteString("OK")
	}
}
