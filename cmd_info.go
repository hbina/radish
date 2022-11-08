package redis

func InfoCommand(c *Client, args [][]byte) {
	c.Conn().WriteBulkString("")
}
