package redis

// https://redis.io/commands/watch/
// WATCH key [key ...]
func WatchCommand(c *Client, args [][]byte) {
	// Currently no-op because we are not multi-threaded to begin with
	c.Conn().WriteString("OK")
}
