package redis

// https://redis.io/commands/multi/
// MULTI
func MultiCommand(c *Client, args [][]byte) {
	// Currently no-op because we are not multi-threaded to begin with
	c.Conn().WriteString("OK")
}
