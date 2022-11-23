package redis

// https://redis.io/commands/watch/
// WATCH key [key ...]
func WatchCommand(c *Client, args [][]byte) {
	if len(args) == 0 {
		c.Conn().WriteError(ZeroArgumentErr)
		return
	}

	// Currently no-op because we are not multi-threaded to begin with
	c.Conn().WriteString("OK")
}
