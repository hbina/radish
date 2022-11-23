package redis

// https://redis.io/commands/debug/
func DebugCommand(c *Client, args [][]byte) {
	c.Conn().WriteString("Not implemented")
}
