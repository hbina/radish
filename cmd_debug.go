package redis

func DebugCommand(c *Client, args [][]byte) {
	c.Conn().WriteString("Not implemented")
}
