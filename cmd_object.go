package redis

import (
	"github.com/tidwall/redcon"
)

func ObjectCommand(c *Client, cmd redcon.Command) {
	c.Conn().WriteNull()
}
