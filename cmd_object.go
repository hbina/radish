package redis

// https://redis.io/commands/object/
func ObjectCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError(ZeroArgumentErr)
		return
	}

	c.Conn().WriteNull()
}
