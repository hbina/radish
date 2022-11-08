package redis

func FunctionCommand(c *Client, args [][]byte) {
	c.Conn().WriteString("OK")
}
