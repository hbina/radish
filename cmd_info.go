package redis

import (
	"github.com/tidwall/redcon"
)

func InfoCommand(c *Client, cmd redcon.Command) {
	c.Conn().WriteBulkString("")
}
