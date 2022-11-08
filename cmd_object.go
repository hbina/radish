package redis

func ObjectCommand(c *Client, args [][]byte) {
	c.Conn().WriteNull()
}
