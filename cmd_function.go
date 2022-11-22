package redis

// https://redis.io/commands/function/
func FunctionCommand(c *Client, args [][]byte) {
	c.Conn().WriteString("Not implemented")
}
