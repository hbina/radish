package redis

import (
	"github.com/redis-go/redcon"
)

func InfoCommand(c *Client, cmd redcon.Command) {
	c.Conn().WriteBulkString("")
}
