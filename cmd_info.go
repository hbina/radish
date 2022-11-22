package redis

// https://redis.io/commands/info/
func InfoCommand(c *Client, args [][]byte) {
	c.Conn().WriteBulkString("OK")
}
